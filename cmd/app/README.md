# go-api-template

Blueprint for Go applications

## 📌 Usage

```sh
MODUL_NAME [-logLevel trace|debug|info|warning|error] [-LogDestination stdout|stderr|null|/path/to/logfile] [-version] [-about] [-help] [-crypt <text>]
```

### 🛠 Available Flags

| **Flag**                   | **Description**                                                |
|----------------------------|----------------------------------------------------------------|
| `-version`                 | Prints the application version and exit                        |
| `-about`                   | Prints details about `MODUL_NAME` and exit                  |
| `-help`                    | Prints this help message and exit                              |
| `-logLevel <level>`        | Set the log level: trace, debug, info, warning ,error          |
| `-logDestination <dest>`   | Set the log destination: stdout, stderr,null, /path/to/logfile |
| `-config </path/file.cfg>` | Specify the path to the config file                            |
| `-crypt <text>`            | Encrypt the given string and exit                              |

---

## 🔍 Examples

### Print Version:

```sh
MODUL_NAME -version
```

### Show About Information:

```sh
MODUL_NAME -about
```

### Enable Debug Mode (Verbose Logging):

```sh
`MODUL_NAME` -logLevel debug -logDestination stdout
```

### Enable Trace Logging (Source Code Location in Logs):

```sh
`MODUL_NAME` -logLevel trace -logDestination stdout
```

### Encrypt a String (`mysecret` in this example):

```sh
`MODUL_NAME` -crypt "mysecret"
🔐 **Output:** Encrypted string (useful for securing credentials).
```

### Get monitoring data from a smart meter:

```sh
curl -k -H "X-Api-Key: 12345678" https://localhost:4000/api/monitoring
```

---

## 📦 Features

✅ **what ever** – explain what ever


---

## 📖 Documentation

For detailed setup and configuration, visit our **[official documentation]**.

---

## 👨‍💻 Contributing

Want to contribute? Feel free to submit **pull requests** or report issues in the repository.

---

## 📜 License

`MODUL_NAME` is licensed under the **MIT License**.

---

## **🌐 IP Address / IP Network Filter**

`MODUL_NAME` allows **IP-based access control** via the configuration file.

- **`blockedIPs`**: Defines **blocked** IP addresses/networks.
- **`allowedIPs`**: Defines **allowlisted** IP addresses/networks.  
  If set to an **empty list** or `"ALL"`, all IP addresses/networks are allowed.

🔹 **Priority Rule:** `blockedIPs` **takes precedence** over `allowedIPs`.

