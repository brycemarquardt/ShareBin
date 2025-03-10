/*
This file is part of ShareBin.

ShareBin is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

ShareBin is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with ShareBin. If not, see <https://www.gnu.org/licenses/>.
*/

// Elements in main page
let fileInput;
let customFileUpload;
let password;
let duration;
let durationModifiers;
let burn;
let textArea;
let uploadButton;

let files;

// Elements in upload progress page
let progressPage;
let uploadPercent;
let uploadInfo;
let link;

let convertMinutes = {
    "minutes": 1,
    "hours": 60,
    "days": 1440,
    "months": 43800,
    "years": 525960
}

window.onbeforeunload = () => {
    fetch('/deleteSession', { method: 'POST' })
}

window.onload = () => {
    // Elements in main page
    fileInput = document.getElementById('fileInput');
    customFileUpload = document.getElementById('customUpload');
    password = document.getElementById("password");
    burn = document.getElementById("burn");
    duration = document.getElementById("duration");
    durationModifiers = document.getElementById("durationModifiers");
    textArea = document.getElementById("textarea");
    uploadButton = document.getElementById("upload")

    // Elements in upload progress page
    progressPage = document.getElementById("progressPage");
    uploadPercent = document.getElementById("uploadPercent");
    uploadInfo = document.getElementById("uploadInfo");
    link = document.getElementById("link");

    // Load saved preferences or set defaults
    duration.value = localStorage.getItem("duration") || 5;
    durationModifiers.value = localStorage.getItem("durationModifiers") || "minutes";

    duration.addEventListener('input', (e) => {
        localStorage.setItem('duration', e.target.value);
    });

    durationModifiers.addEventListener('change', (e) => {
        localStorage.setItem('durationModifiers', e.target.value);
    });

    textArea.addEventListener('input', () => {
        files = null;
        customFileUpload.innerHTML = "Click here, drag & drop, or Ctrl + V anywhere to select file"
        uploadButton.innerHTML = "Upload text (Ctrl + Enter)";
    });

    textArea.addEventListener('keydown', (e) => {
        if (e.key === 'Tab') {
            e.preventDefault();
            let start = textArea.selectionStart;
            let end = textArea.selectionEnd;
            textArea.value = textArea.value.substring(0, start) + "\t" + textArea.value.substring(end);
            textArea.selectionStart = textArea.selectionEnd = start + 1;
        }
    });

    textArea.addEventListener('paste', (event) => {
        event.preventDefault();
        const pasteData = (event.clipboardData || window.clipboardData).getData('text');
        textArea.value = textArea.value + pasteData;
        textArea.selectionStart = textArea.selectionEnd = textArea.selectionStart + pasteData.length;
    });

    customFileUpload.addEventListener('click', () => {
        fileInput.click()
    });

    fileInput.addEventListener('change', function () {
        handleFiles(this.files);
    });

    document.addEventListener('drop', function (event) {
        event.preventDefault();
        handleFiles(event.dataTransfer.files);
    });

    document.addEventListener('paste', function (event) {
        event.preventDefault();
        if (event.clipboardData && event.clipboardData.files.length > 0) {
            handleFiles(event.clipboardData.files);
        }
    });

    document.addEventListener('keydown', function (event) {
        if (event.ctrlKey && event.key === 'Enter') {
            upload();
            event.preventDefault();
        }
    });

    textArea.focus()
}

function copyLink() {
    copyToClipboard(link.innerHTML)
}

function handleFiles(f) {
    if (f && f.length > 0) {
        textArea.value = "";
        customFileUpload.innerHTML = f.length + " files selected";
        uploadButton.innerHTML = "Upload file";
        files = f;
    }
}

function upload() {
    if (files || textArea.value.trim() !== "") {
        const xhr = new XMLHttpRequest();

        xhr.onloadstart = function () {
            progressPage.style.display = "flex";
        };

        xhr.upload.addEventListener('progress', function (event) {
            if (event.lengthComputable) {
                const percentComplete = Math.round((event.loaded / event.total) * 100);
                if (percentComplete !== 100) {
                    uploadPercent.innerHTML = percentComplete + "%";
                } else {
                    uploadPercent.innerHTML = percentComplete + "% (waiting for link...)";
                }
                console.log(`Upload progress: ${percentComplete}%`);
            }
        });

        // Handle successful upload
        xhr.addEventListener('load', function () {
            if (xhr.status === 200) {
                uploadPercent.innerHTML = "100%";
                link.innerHTML = xhr.responseText;
                uploadInfo.style.display = "flex";

                // Use server-side QR code generation instead of client-side QRCode
                const shareID = xhr.responseText.split('/').pop(); // Extract share ID from URL
                qrCode.innerHTML = `<img src="/qr/${shareID}" alt="QR Code for ${shareID}" style="max-width: 128px;">`;
            } else {
                console.log('Upload failed with status:', xhr.status);
            }
        });

        // Handle upload errors
        xhr.addEventListener('error', function () {
            document.getElementById("uploadPercent").innerHTML = "Upload error";
        });

        // Handle aborts
        xhr.addEventListener('abort', function () {
            console.log('Upload aborted');
        });

        let minutes = duration.value * convertMinutes[durationModifiers.value]

        if (files) {
            const formData = new FormData();
            formData.append("duration", minutes);
            formData.append("pass", password.value);
            formData.append("burn", burn.checked);

            // Append all files to the FormData object
            for (let i = 0; i < files.length; i++) {
                formData.append('file', files[i]); // 'file' as the key
            }

            xhr.open('POST', '/');
            xhr.send(formData);
        } else {
            if (textArea.value.trim() !== "") {
                xhr.open('POST', '/postText');
                xhr.send(JSON.stringify({ duration: minutes, pass: password.value, burn: burn.checked, text: textArea.value }));
            }
        }
    }
}
