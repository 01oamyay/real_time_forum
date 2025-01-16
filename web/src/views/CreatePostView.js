import AbstractView from "./AbstractView.js";
import fetcher from "../pkg/fetcher.js";
import redirect from "../index.js";

const path = "/api/post/create";

const createPost = async (title, text, categories) => {
  let body = {
    title: title,
    data: text,
    categories: categories,
  };

  const data = await fetcher.post(path, body);
  if (data && data.msg !== undefined) {
    let showErr = document.getElementById("showError");
    showErr.innerHTML = data.msg;
    return;
  }
  redirect.navigateTo(`/post/${data.post_id}`);
};

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Create-post");
  }

  async getHtml() {
    const categories = await this.getCategories();

    return `
        <div class="container">
        <form id="form-createPost" class="form-createPost <text-center>" onsubmit="return false;">
        <div class="mb-3">
            <label for="TitleInput" class="form-label">Title</label>
            <input maxlength="58" name="title" type="text" class="form-control" id="TitleInput" required>
            <div class="form-text">Maximum of 58 characters</div>
        </div>
        <div class="mb-3">
            <label for="TextInput" class="form-label">Text</label>
            <textarea name="text" maxlength="10000" id="TextInput" rows="3" required></textarea>
            <div class="form-text">Maximum of 10000 characters</div>
        </div>
        <div class="mb-3">
        <label>Select Categories:</label>
           <div class="multi-select">
                ${categories
                  .map(
                    (category) =>
                      ` <input type="checkbox" id="${category.name}" name="${category.name}" value="${category.name}">
                        <label for="${category.name}">${category.name}</label><br>`
                  )
                  .join("")}
            </div>
        </div>
        <button class="btn btn-primary">Post</button>
        <div id="showError" class="error-message"></div>
        </form>
        </div>
        `;
  }

  async init() {
    const createPostForm = document.getElementById("form-createPost");
    createPostForm.addEventListener("submit", function (e) {
      e.preventDefault();

      const title = document.getElementById("TitleInput").value.trim();
      const data = document.getElementById("TextInput").value.trim();

      const getCheckedCategories = () => {
        const checkedCheckboxes = document.querySelectorAll(
          '.multi-select input[type="checkbox"]:checked'
        );
        return Array.from(checkedCheckboxes).map((checkbox) => checkbox?.value);
      };

      const categories = getCheckedCategories();

      // Validate inputs
      if (!title || !data || !categories?.length) {
        document.getElementById("showError").textContent =
          "Please fill in all fields";
        return;
      }

      createPost(title, data, categories);
    });
  }
}
