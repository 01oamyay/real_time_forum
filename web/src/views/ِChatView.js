import AbstractView from "./AbstractView.js";

let limit = 100;
let offset = 0;
let ended = false;

async function GetMessages(userID) {
  const res = await fetch(
    `/api/chat/${userID}?limit=${limit}&offset=${offset}`
  );
  console.log(res.status);
  const messages = await res.json();
  console.log(messages);

  if (messages?.status == 400 && offset > 0) {
    ended = true;
    return [];
  }

  if (messages?.msg) {
    return;
  }

  if (!messages?.length && offset > 0) {
    ended = true;
    return [];
  }

  offset += messages?.length;

  return messages;
}

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Chat");
  }

  async getHtml() {
    return `<div id="chat">
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
    const receiver_id = this.params.userID;
    const messages = await GetMessages(receiver_id);

    let sender_id;
    if (receiver_id == messages.chat.user_id) {
      sender_id = messages.chat.user_id_1;
    } else {
      sender_id = messages.chat.user_id;
    }

    const sendBtn = document.getElementById("send");
    const input = document.querySelector(".chat__conversation-panel__input");
    input.addEventListener("keydown", (e) => {
      if (e.key === "Enter") {
        sendBtn.click();
      }
    });
    sendBtn.addEventListener("click", async () => {
      let message = document.querySelector(".chat__conversation-panel__input");
      if (message?.value) {
        const sendEvent = new CustomEvent("send-msg", {
          detail: {
            chat_id: messages?.chat?.id,
            sender_id: sender_id,
            content: message?.value,
            created_at: new Date().toISOString(),
          },
        });
        document.dispatchEvent(sendEvent);
        message.value = "";
      }
    });

    document.addEventListener("msg", (e) => {
      if (e.detail.chat_id == messages.chat.id) {
        insertMsg([e.detail.data], sender_id);
      }
    });

    insertMsg(messages.messages, sender_id);
  }
}

function insertMsg(messages, sender_id) {
  let chatContainer = document.querySelector(".chat__conversation-board");
  if (messages.length == 0) {
    chatContainer.innerHTML = `<p class="noMsg">No Messag</p>`;
  } else {
    chatContainer.innerHtml = "";
    messages.forEach((message) => {
      const msg = `
        <div class="chat__conversation-board__message-container ${
          message.sender_id == sender_id ? "reversed" : ""
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
              <span>${message.content}</span>
            </div>
          </div>
        </div>
      `;

      chatContainer.appendChild(
        document.createRange().createContextualFragment(msg)
      );
    });
  }
}
