import fetcher from "../pkg/fetcher.js";


export default class {
  constructor(params) {
    this.params = params;
    this.categories = []; // Store categories
  }

  setTitle(title) {
    document.title = title;
  }

  addStyle(fileName) {
    // Check if style already exists
    if (document.querySelector(`link[data-view-style="${fileName}"]`)) {
      return;
    }

    const link = document.createElement("link");
    link.rel = "stylesheet";
    link.setAttribute("data-view-style", fileName);

    const filePath = `/src/assets/css/${fileName}.css`;
    link.href = filePath;

    document.head.appendChild(link);
  }

  removeStyles() {
    const viewStyles = document.querySelectorAll("link[data-view-style]");
    viewStyles.forEach((style) => style.remove());
  }

  async getCategories() {
    // If categories are already fetched, return them
    if (this.categories && this.categories.length > 0) {
      return this.categories;
    }

    try {
      const fetchedCategories = await fetcher.get("/api/categories");

      // Validate the fetched categories
      if (Array.isArray(fetchedCategories)) {
        this.categories = fetchedCategories.filter(
          (cat) => cat?.name != "General"
        );
        return this.categories;
      } else {
        console.error("Fetched categories is not an array:", fetchedCategories);
        return [];
      }
    } catch (error) {
      console.error("Error fetching categories:", error);
      return [];
    }
  }

  async getHtml() {
    return "";
  }

  async init() {}
}
