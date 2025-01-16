import AbstractView from "./AbstractView.js";

export default class extends AbstractView {
  constructor(params, user) {
    super(params);
    this.user = user;
  }

  async getHtml() {
    const categories = (await this.getCategories()) || [];
    const isAuthorized = Boolean(this.user.id);

    let sidebarHtml = `
      <div id="sidebar" class="sidebar">
        ${
          isAuthorized
            ? `
          <div class="sidebar-section">
            <div class="sidebar-section-title">My Posts</div>
            <div class="sidebar-menu-item" data-category="my-posts">My Posts</div>
            <div class="sidebar-menu-item" data-category="liked-posts">Liked Posts</div>
            <div class="sidebar-menu-item" data-category="disliked-posts">Disliked Posts</div>
          </div>
        `
            : ""
        }

        <div class="sidebar-section">
          <div class="sidebar-section-title">Discover Categories</div>
          ${categories
            .map(
              (category) => `
            <div class="sidebar-menu-item" data-category="${category.name}">
              ${category.name}
            </div>
          `
            )
            .join("")}
        </div>
        <div class="togel-sidbar">
          <button id="togel-sidbar" class="togel-sidbar">
          <span> Filter by catigores <img src="/src/assets/img/settings-sliders.svg" alt="filter"></span>
          </button>
        </div>
      </div>
    `;

    return sidebarHtml;
  }

  async init() {
    const sidebarItems = document.querySelectorAll(".sidebar-menu-item");
    sidebarItems.forEach((item) => {
      item.addEventListener("click", (e) => {
        const category = e.target.getAttribute("data-category");

        // For user-specific sections, check authorization
        const isAuthorized = Boolean(this.user.id);

        if (
          !isAuthorized &&
          ["my-posts", "liked-posts", "disliked-posts"].includes(category)
        ) {
          // Redirect to login or show error
          window.location.href = "/sign-in";
          return;
        }

        // Dispatch a custom event that HomeView can listen to
        const event = new CustomEvent("category-selected", {
          detail: { category: category },
        });
        document.dispatchEvent(event);
      });
    });

    // Toggle sidebar
    const check = document.querySelector("#togel-sidbar");
    if (check) {
      check.addEventListener("click", () => {
        const sidebarSection = document.querySelectorAll(".sidebar-section");
        sidebarSection.forEach((item) => {
          switch (item.style.display) {
            case "block":
              item.style.display = "none";
              break;
            case "none":
              item.style.display = "block";
              break;
            default:
              item.style.display = "block";
          }
          // on resize
          window.addEventListener("resize", () => {
            if (window.innerWidth > 768) {
              item.style.display = "block";
            }
          }
          );
        });
      });
    }
  }
}
