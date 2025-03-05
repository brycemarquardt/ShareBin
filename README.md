# ShareBin
File sharing, URL shortener, and pastebin all in one place with QR code and curl support. Uses stream-based cryptography and data processing that can handle gigabytes of data with fixed memory and CPU usage. It can run on any platform, including PaaS like Replit or Render, and is highly customizable.  
<ins>**Please star this project if you find it useful, thank you!**</ins>

| Feature                          | Description                                      |
|----------------------------------|--------------------------------------------------|
| Mobile-friendly                  | Works in mobile browsers, supports file/text upload via Ctrl+V, drag and drop, file browsing, or terminal |
| Authentication                   | Supports password authentication for both upload and download |
| Easy Setup                       | Requires only `go build .` or use the `docker-compose.yaml` to get started |
| Customizable                     | Easily modify styles by replacing `static/theme.css` with a file from [classless CSS](https://github.com/dbohdan/classless-css) or adjust the well-commented HTML layout |
| Cross-Platform                   | Runs on any OS or deployment platforms like Replit, Render, Fly.io, etc. |
| Secure Encryption                | Password-protected data secured with AES and PBKDF2 |
| On-the-Fly Decryption            | Encrypted data is decrypted in memory, never written to disk |
| Short, Unambiguous URLs          | Generates concise, collision-free URLs (omitting ambiguous characters like i, I, l, 1) |
| QR Code Support                  | Quick sharing of files and URLs to/from mobile devices via QR codes |
| Dashboard                        | View, filter, sort, and manage uploaded shares with a web-based dashboard |

# Why ShareBin is Better :zap:
ShareBin combines file sharing, URL shortening, and pastebin functionality with advanced features like QR code generation, secure encryption, and a user-friendly dashboard, making it ideal for personal and professional use.

# URL Shortener 🔗
Simply paste any valid URL (must start with `http://` or `https://`) into the text box and upload to create a short, shareable link.

# Don’t Like How It Looks? 🎨
Pick a `.css` file from [classless CSS](https://github.com/dbohdan/classless-css) and replace `static/theme.css`, or search for "classless CSS" to find other options. The HTML pages are well-commented and structured for easy customization.

# How to Build with Docker :whale2:
1. Clone or download this repository
2. Create a folder named `uploads`
3. Run `docker compose up`

# How to Build Without Docker 📟
1. Clone or download this repository
2. Open a terminal in the project directory
3. Run `go build .`

# Settings ⚙️
You can modify the variables inside `data/settings.json`:
- `FileSizeLimitMB`: Limit file size (in megabytes)
- `TextSizeLimitMB`: Limit text size (in megabytes)
- `StreamSizeLimitKB`: Limit file encryption, decryption, upload, and download buffer stream size (in KB) to control memory usage
- `StreamThrottleMS`: Add a throttle to encryption, decryption, upload, and download buffers to limit CPU usage (in milliseconds)
- `Pbkdf2Iterations`: Key derivation algorithm iterations; higher values increase security (e.g., 100,000 is recommended)
- `CmdUploadDefaultDurationMinute`: Default file duration (in minutes) for curl uploads if duration isn’t specified
- `enablePassword`: Enable or disable password protection for site authentication
- `password`: Password value for site authentication; use a strong, long password or consider external authentication for enhanced security

You can tune CPU/memory usage by calculating memory usage per second with `streamSizeLimitKB * (1000/streamThrottleMS)`. The default settings can handle 40 MB of data/second for file upload, download, encryption, and decryption—adjust as needed.

# Curl Upload ⬆️
Example: `curl -F file=@main.go -F duration=10 -F pass=123 -F burn=true https://yoursite.com`  
Note: The `duration`, `pass`, and `burn` parameters are optional. For a quick upload, use `curl -F file=@file.txt https://yoursite.com`. If your site is password-protected, add `-F auth=yourpassword`.

# Security 🔒
For maximum security, encrypt your files before uploading. ShareBin uses AES and PBKDF2 for password-protected data, with on-the-fly decryption to prevent decrypted data from being written to disk.

# Contribution 🤝
Feel free to open an issue if you have a feature idea or send a pull request. Your contributions are welcome!

# License
This project is licensed under the GNU General Public License (GPL) v3 or later. See [LICENSE](LICENSE) for details.
