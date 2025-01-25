import fetcher from "../pkg/fetcher.js";
import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params, user) {
    super(params);
    this.user = user;
    this.contacts = [];
  }

  async getContacts() {
    const contacts = await fetcher.get("/api/contacts");

    contacts.sort((a, b) => {
      // Prioritize contacts with valid last messages
      if (a.has_last_msg !== b.has_last_msg) {
        return b.has_last_msg - a.has_last_msg;
      }

      // If both have last messages, sort by time (most recent first)
      if (a.has_last_msg) {
        return new Date(b.last_msg.Time) - new Date(a.last_msg.Time);
      }

      // If no last messages, sort alphabetically by firstName
      return a.firstName.localeCompare(b.firstName);
    });
    this.contacts = contacts;
  }

  async getHtml() {
    let sidebarHtml = `
        <div id="user_sidebar" class="user_sidebar">
            <ul>
        <div class="user_sidebar-section-title">Users</div>
    `;

    await this.getContacts();

    console.log(this.contacts);

    this.contacts.forEach((user) => {
      sidebarHtml += `
            <li class="item">
                <a href="/chat/${user.user_id}" data-link>
                    <div class="item-content">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                        </svg>
                        <span>${(user.firstName, user.lastName)}</span>
                    </div>
                    <span id=${user.user_id} class="dot ${
        user.isOnline ? "online" : "offline"
      }"></span>
                </a>
            </li>
        `;
    });

    sidebarHtml += `
            </ul>
        </div>
        <div class="toggle-user_sidebar">
            <button id="toggle-user_btn" class="toggle-user_btn">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960"><path d="M40-160v-112q0-34 17.5-62.5T104-378q62-31 126-46.5T360-440q66 0 130 15.5T616-378q29 15 46.5 43.5T680-272v112H40Zm720 0v-120q0-44-24.5-84.5T666-434q51 6 96 20.5t84 35.5q36 20 55 44.5t19 53.5v120H760ZM360-480q-66 0-113-47t-47-113q0-66 47-113t113-47q66 0 113 47t47 113q0 66-47 113t-113 47Zm400-160q0 66-47 113t-113 47q-11 0-28-2.5t-28-5.5q27-32 41.5-71t14.5-81q0-42-14.5-81T544-792q14-5 28-6.5t28-1.5q66 0 113 47t47 113ZM120-240h480v-32q0-11-5.5-20T580-306q-54-27-109-40.5T360-360q-56 0-111 13.5T140-306q-9 5-14.5 14t-5.5 20v32Zm240-320q33 0 56.5-23.5T440-640q0-33-23.5-56.5T360-720q-33 0-56.5 23.5T280-640q0 33 23.5 56.5T360-560Zm0 320Zm0-400Z"/></svg>
            </button>
        </div>
    `;

    return sidebarHtml;
  }

  async init() {
    const check = document.querySelector("#toggle-user_btn");
    if (check) {
      const sidebar = document.querySelector("#user_sidebar");
      check.addEventListener("click", () => {
        sidebar.style.display =
          sidebar.style.display === "flex" ? "none" : "flex";
      });
      window.addEventListener("resize", () => {
        if (window.innerWidth < 768) {
          sidebar.style.display = "none";
        } else {
          sidebar.style.display = "flex";
        }
      });
    }

    document.addEventListener("user-online", (e) => {
      console.log("user list:", e.detail);
      const span = document.getElementById(`${e.detail}`);
      span?.classList.remove("offline");
      span.classList.add("online");
    });

    document.addEventListener("user-offline", (e) => {
      const span = document.getElementById(`${e.detail}`);
      span?.classList.remove("online");
      span.classList.add("offline");
    });
  }
}
