import Utils from "../pkg/Utils.js";
import AbstractView from "./AbstractView.js";

let limit = 10;
let offset = 0;
let ended = false;
let sender_id;
let receiver_id;

async function GetMessages() {
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

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Chat");
    this.typingTimeout = null;
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
        <textarea class="chat__conversation-panel__input panel-item" placeholder="Type a message..." rows="1"></textarea>
        <button id="send" class="chat__conversation-panel__button panel-item btn-icon send-message-button">    
      <svg xmlns="http://www.w3.org/2000/svg" height="24px" viewBox="0 -960 960 960" width="24px"><path d="M120-160v-640l760 320-760 320Zm80-120 474-200-474-200v140l240 60-240 60v140Zm0 0v-400 400Z"/></svg>
        </button>
    </div>
  `;
  }

  async init() {
    offset = 0;
    receiver_id = this.params.userID;
    const messages = await GetMessages(receiver_id);
    if (!messages) return;

    this.chatID = messages?.chat?.id;

    const inputField = document.querySelector(
      "textarea.chat__conversation-panel__input.panel-item"
    );

    if (receiver_id == messages.chat.user_id) {
      sender_id = messages.chat.user_id_1;
    } else {
      sender_id = messages.chat.user_id;
    }

    inputField.setAttribute("data-sender_id", `${sender_id}`);

    const sendBtn = document.getElementById("send");
    const msgInput = document.querySelector(
      `.chat__conversation-panel__input[data-sender_id="${sender_id}"]`
    );
    let pendingMessage = "";

    document.addEventListener("keydown", (e) => {
      if (e.key === "Enter" && !e.shiftKey) {
        sendBtn.click();
        const typingMsg = document.getElementById("loading");
        if (typingMsg) {
          typingMsg.remove();
        }
      }
    });

    sendBtn.addEventListener("click", () => {
      if (msgInput?.value && msgInput.dataset.sender_id == sender_id) {
        pendingMessage = msgInput.value;
        const sendEvent = new CustomEvent("send-msg", {
          detail: {
            event: "msg",
            payload: {
              chat_id: messages?.chat?.id,
              sender_id: sender_id,
              content: pendingMessage,
              created_at: new Date().toISOString(),
            },
          },
        });
        msgInput.value = "";
        document.dispatchEvent(sendEvent);
        const typingMsg = document.getElementById("loading");
        if (typingMsg) {
          typingMsg.remove();
        }
      }
    });

    inputField.addEventListener("input", (e) => {
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
            chat_id: this.chatID,
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
              chat_id: this.chatID,
              is_typing: false,
              user_id: sender_id,
            },
          },
        });

        document.dispatchEvent(typingEvent);
      }, 1000);
    });

    document.addEventListener("typing", (e) => {
      const status = e.detail;

      let loadingElem = document.getElementById("loading");

      if (status.chat_id != this.chatID) return;

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
        <span class="nickname">"to fix"</span>
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
    });

    document.addEventListener("msg", (e) => {
      if (e.detail.chat_id == messages.chat.id) {
        const typingMsg = document.getElementById("loading");
        if (typingMsg) {
          typingMsg.remove();
        }
        insertMsg(e.detail, sender_id);
      }
    });

    document.addEventListener("msg-error", (e) => {
      Utils.showToast(e.detail?.error);
      document.querySelector(".chat__conversation-panel__input").value =
        pendingMessage;
    });

    insertMsg(messages.messages, sender_id);

    const chatBoard = document.querySelector(".chat__conversation-board");
    // setup scroll event
    chatBoard.addEventListener("scroll", handleScroll);
  }
}

let isThrottled = false;

async function handleScroll(e) {
  if (ended) {
    const chatBoard = document.querySelector(".chat__conversation-board");
    chatBoard.removeEventListener("scroll", handleScroll);
    return;
  }
  if (isThrottled) return;
  let container = e.target;

  if (container.scrollTop <= 0) {
    isThrottled = true;

    const height = container.scrollHeight;

    try {
      const messages = await GetMessages(receiver_id);
      if (messages?.length) {
        return;
      }

      insertMsg(messages, sender_id, true);

      const newHeight = container.scrollHeight;
      const diff = newHeight - height;
      container.scrollTop = diff;
    } catch (error) {
      console.error("Error loading messages:", error);
    } finally {
      setTimeout(() => {
        isThrottled = false;
        if (container.scrollTop <= 0) {
          handleScroll({ target: container });
        }
      }, 1000);
    }
  }
}

function insertMsg(message, sender_id, pre = false) {
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

function formatDate(dateString) {
  const date = new Date(dateString);
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");
  return `${hours}:${minutes}`;
}
