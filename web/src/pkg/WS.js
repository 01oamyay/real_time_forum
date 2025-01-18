export default class {
  constructor() {
    this.ws = null;
    window.addEventListener("beforeunload", () => {
      if (this.ws && this.ws.readyState == WebSocket.OPEN) {
        this.ws.close(1000, "user leaving");
      }
    });

    document.getElementById("logout-btn")?.addEventListener("click", () => {
      if (this.ws && this.ws.readyState == WebSocket.OPEN) {
        this.ws.close(1000, "user leaving");
      }
    });
  }

  init() {
    this.ws = new WebSocket("/ws");

    console.log("ws enabled");

    this.ws.onclose = (e) => {
      console.log("connection closed");
    };

    this.ws.onerror = (e) => {
      console.log("error in ws", e);
    };

    this.ws.onmessage = (e) => {
      let data = JSON.parse(e.data);
      console.log(data);
      let onlineEvent = new CustomEvent(data.event, {
        detail: data,
      });
      document.dispatchEvent(onlineEvent);
    };

    window.addEventListener("beforeunload", () => {
      if (this.ws && this.ws.readyState == WebSocket.OPEN) {
        this.ws.close(1000, "user leaving");
      }
    });
  }
}
