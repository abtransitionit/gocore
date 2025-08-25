# Todo
* **Purpose:** Provides a **cross-platform way** to retrieve basic system properties (Linux, Windows, macOS).
* **Main functions:**

  * `GetPropertyLocal(property string, params ...string)`: fetches a property from the **core set** (`coreProperties` map).
  * `coreProperties` map contains handlers for properties like:

  | Property       | Description                     |
  |----------------|---------------------------------|
  | `ostype`       | OS type (e.g., linux, windows, darwin) |
  | `osarch`       | OS architecture (e.g., amd64, arm64)  |
  | `cpu`          | Number of CPU cores             |
  | `path`         | PATH environment variable       |
  | `osversion`    | OS version                      |
  | `oskversion`   | OS kernel version               |
  | `osuser`       | Current user                    |
  | `ram`          | Total system RAM in GB          |

* **Notes:**

  * Uses `gopsutil` for system info.
  * Can be used on **any OS**, hence cross-platform.
  * cf. the same library in the [golinux](https://github.com/abtransitionit/golinux/) project. 



### **✅ Key Points**

* Both libraries have the **same design pattern**:

  * `GetPropertyLocal` function.
  * A `map[string]PropertyHandler]` mapping property names to handler functions.
* `gocore` is **generic and cross-platform**.
* `golinux` is **Linux-specific**, providing extra details that `gocore` cannot.
* The separation allows code to:

  * Query common properties with `gocore`.
  * Query Linux-specific properties only when running on Linux, without breaking cross-platform compatibility.

---

If you want, I can **draw a simple diagram showing how gocore and golinux interact** in a system property retrieval workflow. It might make it even clearer. Do you want me to do that?


a pure GO function to get the OS type:
  - darwin
  - linux
  - windows

