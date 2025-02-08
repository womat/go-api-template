# go-api-template
Blueprint for Go applications

## 📌 Usage

```sh
s0counter [-debug] [-trace] [-version] [-about] [-help] [-crypt <text>]
```

### 🛠 Available Flags

| **Flag**        | **Description**                                    |
|-----------------|----------------------------------------------------|
| `-version`      | Prints the application version and exits           |
| `-about`        | Displays details about `s0counter` and exits       |
| `-debug`        | Enables verbose debug logging to stdout            |
| `-trace`        | Enables source code location logging for debugging |
| `-help`         | Prints this help message                           |
| `-crypt <text>` | Encrypts the given string and exits                |

---

## 🔍 Examples

### Print Version:

```sh
s0counter -version
```

### Show About Information:

```sh
s0counter -about
```

### Enable Debug Mode (Verbose Logging):

```sh
s0counter -debug
```

### Enable Trace Logging (Source Code Location in Logs):

```sh
s0counter -trace
```

### Encrypt a String (`mysecret` in this example):

```sh
s0counter -crypt "mysecret"
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

For detailed setup and configuration, visit our **[official documentation](#)**.

---

## 👨‍💻 Contributing

Want to contribute? Feel free to submit **pull requests** or report issues in the repository.

---

## 📜 License

`s0counter` is licensed under the **MIT License**.

---

## **🌐 IP Address / IP Network Filter**

<MODULE> allows **IP-based access control** via the configuration file.

- **`blockedIPs`**: Defines **blocked** IP addresses/networks.
- **`allowedIPs`**: Defines **allowlisted** IP addresses/networks.  
  If set to an **empty list** or `"ALL"`, all IP addresses/networks are allowed.

🔹 **Priority Rule:** `blockedIPs` **takes precedence** over `allowedIPs`.

