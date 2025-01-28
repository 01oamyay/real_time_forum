import Utils from "../pkg/Utils.js";
import AbstractView from "./AbstractView.js";

let limit = 10;
let offset = 0;
let ended = false;
let sender_id;
let receiver_id;
let chatID;

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Chat");
    this.typingTimeout = null;
    this.listeners = false;
    this.sendBtn = null;
    this.msgInput = null;
    this.chat = null;
    this.inputField = null;
    this.handleScroll = this.handleScroll.bind(this);

    // Bind methods
    this.keyDownHandler = this.keyDownHandler.bind(this);
    this.sendMsg = this.sendMsg.bind(this);
    this.inputEvent = this.inputEvent.bind(this);
    this.typingEvent = this.typingEvent.bind(this);
    this.msgEvent = this.msgEvent.bind(this);
    this.msgError = this.msgError.bind(this);
  }

  async getMessages() {
    const res = await fetch(
      `/api/chat/${receiver_id}?limit=${limit}&offset=${offset}`
    );

    const messages = await res.json();

    if (messages.status == 404) {
      Utils.showError(messages.status, messages.msg);
    }

    if (messages?.status == 400 && offset > 0) {
      return [];
    }

    if (messages?.msg) {
      return;
    }

    if (!messages?.messages?.length && offset > 0) {
      ended = true;
      return [];
    }

    offset += messages?.messages?.length || 0;

    messages?.messages?.sort((a, b) => {
      return new Date(a.created_at) - new Date(b.created_at);
    });
    return messages;
  }

  async getHtml() {
    return `  
    <div id="chat">
      <div id="toast">
      <div id="toast-message"></div>
      </div>
      <div class="chat__conversation-board">
      </div>
      <div class="chat__conversation-panel">
        <textarea class="chat__conversation-panel__input panel-item" placeholder="Type a message..." rows="1" max-length="200" required></textarea>
        <button id="send" class="chat__conversation-panel__button panel-item btn-icon send-message-button">    
      <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px"><path d="M120-160v-640l760 320-760 320Zm80-120 474-200-474-200v140l240 60-240 60v140Zm0 0v-400 400Z"/></svg>
        </button>
    </div>
  `;
  }

  destroy() {
    // Remove all event listeners
    document.removeEventListener("keydown", this.keyDownHandler);
    this.sendBtn?.removeEventListener("click", this.sendMsg);
    this.inputField?.removeEventListener("input", this.inputEvent);
    document.removeEventListener("typing", this.typingEvent);
    document.removeEventListener("msg", this.msgEvent);
    document.removeEventListener("msg-error", this.msgError);

    const chatBoard = document.querySelector(".chat__conversation-board");
    chatBoard?.removeEventListener("scroll", this.handleScroll);

    // Clear timeouts
    if (this.typingTimeout) {
      clearTimeout(this.typingTimeout);
    }

    // Reset all variables
    this.listeners = false;
    this.sendBtn = null;
    this.msgInput = null;
    this.chat = null;
    this.inputField = null;
    this.typingTimeout = null;

    // Reset global variables
    offset = 0;
    ended = false;
    chatID = null;
  }

  insertMsg(message, sender_id, pre = false) {
    let chatContainer = document.querySelector(".chat__conversation-board");
    function createMessageHTML(msg) {
      return `
        <div class="chat__conversation-board__message-container ${
          msg.sender_id == sender_id ? "reversed" : ""
        }">
          <div class="chat__conversation-board__message__person">
            <div class="chat__conversation-board__message__person__avatar">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
              </svg>
            </div>
          </div>
          <div class="chat__conversation-board__message__context">
            <div class="chat__conversation-board__message__bubble">
              <div class="chat__conversation-board__message__person__info">
                <span class="nickname">${msg.nickname} - </span>
                <span class="created-at">${formatDate(msg.created_at)}</span>
              </div>
              <p class="message">${msg.content.trim()}</p>
            </div>
          </div>
        </div>
      `;
    }

    if (Array.isArray(message)) {
      if (message.length == 0 && offset == 0) {
        chatContainer.innerHTML = `<p class="noMsg">No Messages</p>`;
        return;
      }

      if (!pre) {
        chatContainer.innerHTML = "";
      }

      message.forEach((msg) => {
        const msgHtml = createMessageHTML(msg);
        if (pre) {
          chatContainer.prepend(
            document.createRange().createContextualFragment(msgHtml)
          );
        } else {
          chatContainer.appendChild(
            document.createRange().createContextualFragment(msgHtml)
          );
        }
      });

      if (!pre) {
        chatContainer.scrollTop = chatContainer.scrollHeight;
      }
    } else {
      const msgHtml = createMessageHTML(message);
      chatContainer.appendChild(
        document.createRange().createContextualFragment(msgHtml)
      );
      chatContainer.scrollTop = chatContainer.scrollHeight;
    }
  }

  keyDownHandler(e) {
    if (e.key === "Enter" && !e.shiftKey) {
      this.sendMsg();
      const typingMsg = document.getElementById("loading");
      if (typingMsg) {
        typingMsg.remove();
      }
    }
  }

  sendMsg() {
    if (this.msgInput?.value && this.msgInput.dataset.sender_id == sender_id) {
      if (
        this.msgInput.value.length == 0 ||
        this.msgInput.value.trim() == "" ||
        this.msgInput.value.length > 200
      ) {
        Utils.showToast("Message length must be between 1 and 200");
        return;
      }
      const sendEvent = new CustomEvent("send-msg", {
        detail: {
          event: "msg",
          payload: {
            chat_id: this.chat?.id,
            sender_id: sender_id,
            content: this.msgInput.value,
            created_at: new Date().toISOString(),
          },
        },
      });
      this.msgInput.value = "";
      document.dispatchEvent(sendEvent);
      const typingMsg = document.getElementById("loading");
      if (typingMsg) {
        typingMsg.remove();
      }
    }
  }

  inputEvent(e) {
    if (this.typingTimeout) {
      clearTimeout(this.typingTimeout);
    }

    if (e.inputType == "insertLineBreak") {
      return;
    }

    const typingEvent = new CustomEvent("typing", {
      detail: {
        event: "typing",
        payload: {
          chat_id: chatID,
          is_typing: true,
          user_id: sender_id,
        },
      },
    });

    document.dispatchEvent(typingEvent);

    this.typingTimeout = setTimeout(() => {
      const typingEvent = new CustomEvent("typing", {
        detail: {
          event: "typing",
          payload: {
            chat_id: chatID,
            is_typing: false,
            user_id: sender_id,
          },
        },
      });

      document.dispatchEvent(typingEvent);
    }, 1000);
  }

  typingEvent(e) {
    const status = e.detail;
    let loadingElem = document.getElementById("loading");

    if (status.chat_id != chatID) return;

    if (status?.is_typing && !loadingElem) {
      let chatContainer = document.querySelector(".chat__conversation-board");

      const loadignFragement = document.createRange()
        .createContextualFragment(`<div id="loading" class="chat__conversation-board__message-container">
<div class="chat__conversation-board__message__person">
  <div class="chat__conversation-board__message__person__avatar">
    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
      <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
    </svg>
  </div>
</div>
<div class="chat__conversation-board__message__context">
  <div class="chat__conversation-board__message__bubble">
    <div class="chat__conversation-board__message__person__info">
      <span class="nickname">${status.nickname}</span>
    </div>
    <svg version="1.1" id="L4" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" 
      viewBox="0 44 52 12" enable-background="new 0 0 0 0" xml:space="preserve" width="30">
      <circle fill="#fff" stroke="none" cx="6" cy="50" r="6">
        <animate
          attributeName="opacity"
          dur="1s"
          values="0;1;0"
          repeatCount="indefinite"
          begin="0.1"/>    
      </circle>
      <circle fill="#fff" stroke="none" cx="26" cy="50" r="6">
        <animate
          attributeName="opacity"
          dur="1s"
          values="0;1;0"
          repeatCount="indefinite" 
          begin="0.2"/>       
      </circle>
      <circle fill="#fff" stroke="none" cx="46" cy="50" r="6">
        <animate
          attributeName="opacity"
          dur="1s"
          values="0;1;0"
          repeatCount="indefinite" 
          begin="0.3"/>     
      </circle>
    </svg>
    </p>
  </div>
</div>
</div>`);

      chatContainer.append(loadignFragement);
      chatContainer.scrollTop = chatContainer.scrollHeight;
    }

    if (loadingElem && !status?.is_typing) {
      loadingElem.remove();
    }
  }

  msgEvent(e) {
    if (e.detail.chat_id == chatID) {
      const typingMsg = document.getElementById("loading");
      if (typingMsg) {
        typingMsg.remove();
      }
      this.insertMsg(e.detail, sender_id);
    }
  }

  msgError(e) {
    Utils.showToast(e.detail?.error);
  }

  async handleScroll(e) {
    if (ended) {
      const chatBoard = document.querySelector(".chat__conversation-board");
      chatBoard?.removeEventListener("scroll", this.handleScroll);
      return;
    }

    if (this.isThrottled) return;
    let container = e.target;

    if (container.scrollTop <= 0) {
      this.isThrottled = true;

      const height = container.scrollHeight;

      try {
        const messages = await this.getMessages();
        if (!messages || messages?.length === 0) {
          return;
        }

        this.insertMsg(messages, sender_id, true);

        const newHeight = container.scrollHeight;
        const diff = newHeight - height;
        container.scrollTop = diff;
      } catch (error) {
        console.error("Error loading messages:", error);
      } finally {
        setTimeout(() => {
          this.isThrottled = false;
          if (container.scrollTop <= 0) {
            this.handleScroll({ target: container });
          }
        }, 1000);
      }
    }
  }

  setupListeners() {
    if (this.listeners) return;

    document.addEventListener("keydown", this.keyDownHandler);
    this.sendBtn?.addEventListener("click", this.sendMsg);
    this.inputField?.addEventListener("input", this.inputEvent);
    document.addEventListener("typing", this.typingEvent);
    document.addEventListener("msg", this.msgEvent);
    document.addEventListener("msg-error", this.msgError);

    const chatBoard = document.querySelector(".chat__conversation-board");
    chatBoard?.addEventListener("scroll", this.handleScroll);

    this.listeners = true;
  }

  async init() {
    this.destroy();

    offset = 0;
    ended = false;
    receiver_id = this.params.userID;

    const messages = await this.getMessages();
    if (!messages) return;

    this.chat = messages?.chat;
    chatID = messages?.chat?.id;

    this.inputField = document.querySelector(
      "textarea.chat__conversation-panel__input.panel-item"
    );

    if (receiver_id == messages.chat.user_id) {
      sender_id = messages.chat.user_id_1;
    } else {
      sender_id = messages.chat.user_id;
    }

    this.inputField.setAttribute("data-sender_id", `${sender_id}`);

    this.sendBtn = document.getElementById("send");
    this.msgInput = document.querySelector(
      `.chat__conversation-panel__input[data-sender_id="${sender_id}"]`
    );

    this.insertMsg(messages.messages, sender_id);

    this.setupListeners();
  }
}

function formatDate(dateString) {
  const date = new Date(dateString);
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");
  return `${hours}:${minutes}`;
}
