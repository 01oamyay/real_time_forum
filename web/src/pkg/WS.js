import Utils from "../pkg/Utils.js";

export default class {
  constructor() {
    this.ws = null;
    this.listeners = false;
  }

  setupListeners() {
    if (this.listeners) return;
    // Handle clean disconnection
    window.addEventListener("beforeunload", () => {
      if (this.ws) this.ws.close();
    });

    document.addEventListener("send-msg", (e) => {
      this?.ws?.send(JSON.stringify(e.detail));
    });

    document.addEventListener("typing", (e) => {
      this?.ws?.send(JSON.stringify(e.detail));
    });

    this.listeners = true;
  }

  init() {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log("WebSocket already connected");
      return;
    }
    this.ws = new WebSocket("/ws");

    this.setupListeners();

    console.log("ws enabled");

    // Monitor connection status
    this.ws.onclose = (event) => {
      console.log("Connection closed:", event.code, event.reason);
    };

    // Watch for error
    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    this.ws.onmessage = (e) => {
      let data = JSON.parse(e.data);
      switch (data.event) {
        case "new_user":
          let newUserEvent = new CustomEvent("new_user", { detail: data });
          document.dispatchEvent(newUserEvent);
          break;
        case "msg":
          let onlineEvent = new CustomEvent("msg", {
            detail: data.data,
          });

          if (!document.location.pathname.startsWith("/chat/")) {
            Utils.showToast(
              "You received a message from " + data.data.nickname
            );
          }

          document.dispatchEvent(onlineEvent);
          break;
        case "msg-error":
          Utils.showToast(data.error, "error");
          break;
        case "typing":
          let typingEvent = new CustomEvent("typing", {
            detail: data.typing,
          });
          document.dispatchEvent(typingEvent);
          break;
        case "user-online":
          let userOnlineEvent = new CustomEvent("user-online", {
            detail: data.data,
          });
          document.dispatchEvent(userOnlineEvent);
          break;
        case "user-offline":
          let userOfflineEvent = new CustomEvent("user-offline", {
            detail: data.data,
          });
          document.dispatchEvent(userOfflineEvent);
          break;
        case "error":
          console.log(data.error);
          break;
        default:
          console.log("Unknown event:", data.event);
      }
    };
  }
}
