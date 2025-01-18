import AbstractView from "./AbstractView.js";
const users = [
    {
        name: "Mohamed Ali",
        status: "online",
        link: "/users/1"
    },
    {
        name: "Sarah Johnson",
        status: "offline",
        link: "/users/2"
    },
    {
        name: "John Doe",
        status: "online",
        link: "/users/3"
    },
    {
        name: "Jane Smith",
        status: "offline",
        link: "/users/4"
    }
];

export default class extends AbstractView {
    constructor(params, user) {
        super(params);
        this.user = user;
    }

    async getHtml() {

        let sidebarHtml = `
        <div id="user_sidebar" class="user_sidebar">
            <ul>
                <div class="user_sidebar-section-title">Users</div>
    `;

        users.forEach(user => {
            sidebarHtml += `
            <li class="item">
                <a href="${user.link}">
                    <div class="item-content">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
                            <path stroke-linecap="round" stroke-linejoin="round" d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                        </svg>
                        <span>${user.name}</span>
                    </div>
                    <span class="dot ${user.status}"></span>
                </a>
            </li>
        `;
        });

        sidebarHtml += `
            </ul>
            <div class="toggle-user_sidebar">
                <button id="toggle-user_sidebar" class="toggle-user_sidebar">
                    <span> users <img src="/src/assets/img/settings-sliders.svg" alt="filter"></span>
                </button>
            </div>
        </div>
    `;

        return sidebarHtml;
    }

    async init() {
        

        // Toggle sidebar
        const check = document.querySelector("#toggle-user_sidbar");
        if (check) {
            check.addEventListener("click", () => {
                
            });
        }
    }
}
