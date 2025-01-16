import AbstractView from "./AbstractView.js";
import Utils from "../pkg/Utils.js";
import fetcher from "../pkg/fetcher.js";

export default class extends AbstractView {
  constructor(params, user) {
    super(params);
    this.user = user;
  }

  async getHtml() {
    const isAuthorized = await fetcher.isLoggedIn();
    return `
        <div class="logo" id="logo">Forum App</div>


        <div id="nav-auth-links" class="${isAuthorized ? "hidden" : ""}">
            <button id="login-btn" class="btn login">Login</button>
            <button id="register-btn" class="btn register">Register</button>
        </div>
        <div id="nav-user-links" class="${isAuthorized ? "" : "hidden"}">
            ${
              location.pathname !== "/create-post"
                ? '<button id="create-post-btn" class="btn create">Create Post</button>'
                : ""
            }
            <button id="logout-btn" class="btn logout">Logout</button>
        </div>
        `;
  }

  async init() {
    const isAuthorized = await fetcher.isLoggedIn();

    document.getElementById("logo")?.addEventListener("click", () => {
      window.location.href = "/";
    });

    if (isAuthorized) {
      // User-specific buttons
      document
        .getElementById("create-post-btn")
        ?.addEventListener("click", () => {
          window.location.href = "/create-post";
        });

      document.getElementById("logout-btn")?.addEventListener("click", () => {
        Utils.logOut();
      });
    } else {
      // Login/Register buttons
      document.getElementById("login-btn")?.addEventListener("click", () => {
        window.location.href = "/sign-in";
      });

      document.getElementById("register-btn")?.addEventListener("click", () => {
        window.location.href = "/sign-up";
      });
    }
  }
}
