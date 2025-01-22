import Utils from "../pkg/Utils.js";

export default class {
  constructor() {
    this.ws = null;
  }

  init() {
    this.ws = new WebSocket("/ws");

    console.log("ws enabled");

    // Monitor connection status
    this.ws.onclose = (event) => {
      console.log("Connection closed:", event.code, event.reason);
    };

    // Implement reconnection logic if needed
    this.ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    this.ws.onmessage = (e) => {
      let data = JSON.parse(e.data);
      let onlineEvent = new CustomEvent(data.event, {
        detail: data,
      });
      document.dispatchEvent(onlineEvent);

      if (
        data.event == "msg" &&
        !document.location.pathname.startsWith("/chat/")
      ) {
        const sender = data.data.nickname;
        Utils.showToast(`You received a new message from ${sender}`, "msg");
      }
    };
    // Handle clean disconnection
    window.addEventListener("beforeunload", () => {
      if (this.ws) ws.close();
    });

    document.addEventListener("send-msg", (e) => {
      this?.ws?.send(JSON.stringify(e.detail));
    });

    // document.addEventListener();
  }
}
