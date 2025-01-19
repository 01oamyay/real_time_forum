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
      console.log(data);
      let onlineEvent = new CustomEvent(data.event, {
        detail: data,
      });
      document.dispatchEvent(onlineEvent);
    };
    // Handle clean disconnection
    window.addEventListener("beforeunload", () => {
      if (this.ws) ws.close();
    });

    document.addEventListener("send-msg", (e) => {
      console.log("dqzdqz", e.detail);
      this?.ws?.send(JSON.stringify(e.detail));
    });
  }
}
