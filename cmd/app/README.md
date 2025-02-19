# go-api-template

Blueprint for Go applications

## 📌 Usage

```sh
MODUL_NAME [-logLevel debug|info|warning|error] [-LogDestination stdout|stderr|null|/path/to/logfile] [-version] [-about] [-help]
```

### 🛠 Available Flags

| **Flag**                   | **Description**                                                |
|----------------------------|----------------------------------------------------------------|
| `-version`                 | Prints the application version and exit                        |
| `-about`                   | Prints details about `MODUL_NAME` and exit                     |
| `-help`                    | Prints this help message and exit                              |
| `-logLevel <level>`        | Set the log level: debug, info, warning ,error                 |
| `-logDestination <dest>`   | Set the log destination: stdout, stderr,null, /path/to/logfile |
| `-config </path/file.cfg>` | Specify the path to the config file                            |

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

### Enable Debug Logging (Source Code Location in Logs):

```sh
MODUL_NAME -logLevel debug -logDestination stdout
```

### Get monitoring data:

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

## generate a self-signed certificate for development**

    openssl req -x509 -nodes -newkey rsa:2048 -keyout selfsigned.key -out selfsigned.crt -days 35600 -subj "/C=AT/ST=Vienna/L=Vienna/O=ITDesign/OU=DEV/CN=localhost/emailAddress=support@itdesign.at"
      -subj description
       /C=AT								Country
       /ST=Vienna							State (optional).
       /L=Vienna							Location – City (optional).
       /O=company							company (optional).
       /OU=IT								Organizational Unit – (optional).
       /CN=my-domain.com					Common Name – IMPORTANT! your domain name or localhost.
       /emailAddress=admin@my-domain.com	E-Mail-Address (optional).
