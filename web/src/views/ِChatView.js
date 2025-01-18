import AbstractView from "./AbstractView.js";

const newChatElement = () => {
    let postCard = `
<div id="chat">
  <div class="chat__conversation-board">
    <div class="chat__conversation-board__message-container">
      <div class="chat__conversation-board__message__person">
        <div class="chat__conversation-board__message__person__avatar">
          <img src="https://randomuser.me/api/portraits/women/44.jpg" alt="Monika Figi" />
        </div>
      </div>
      <div class="chat__conversation-board__message__context">
        <div class="chat__conversation-board__message__bubble">
          <span>Somewhere stored deep, deep in my memory banks is the phrase "It really whips the llama's ass".</span>
        </div>
      </div>
    </div>


    <div class="chat__conversation-board__message-container reversed">
      <div class="chat__conversation-board__message__person">
        <div class="chat__conversation-board__message__person__avatar">
          <img src="https://randomuser.me/api/portraits/men/9.jpg" alt="Dennis Mikle" />
        </div>
      </div>
      <div class="chat__conversation-board__message__context">
        <div class="chat__conversation-board__message__bubble">
          <span>Winamp's still an essential.</span>
        </div>
      </div>
    </div>
  </div>

  <div class="chat__conversation-panel">
    <div class="chat__conversation-panel__container">
      <input class="chat__conversation-panel__input panel-item" placeholder="Type a message..." />
      <button class="chat__conversation-panel__button panel-item btn-icon send-message-button">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="22" y1="2" x2="11" y2="13"></line>
          <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
        </svg>
      </button>
    </div>
  </div>
</div>

`
    return document.createRange().createContextualFragment(postCard);
};

export default class extends AbstractView {

}
