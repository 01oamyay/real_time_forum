import AbstractView from "./AbstractView.js";
import redirect from "../index.js";
import fetcher from "../pkg/fetcher.js";
import Utils from "../pkg/Utils.js";

const signIn = async (login, password) => {
  let body = {
    login: login,
    password: password,
  };

  const data = await fetcher.post("/api/signin", body);
  if (data && data.msg !== undefined) {
    let showErr = document.getElementById("showError");
    showErr.innerHTML = data.msg;
    return;
  }

  const payload = Utils.parseJwt(data.token);
  localStorage.setItem("id", payload.id);
  localStorage.setItem("role", redirect.roles.user);
  redirect.navigateTo("/");
};

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Sign In");
  }

  async getHtml() {
    return `
            <div class="auth-wrapper">
                <form id="sign-in-form" class="auth-form">
                    <h2 class="form-title">Sign In</h2>
                    
                    <div class="form-group">
                        <label for="login">Email / Nickname</label>
                        <input 
                            type="text" 
                            id="login" 
                            class="form-control" 
                            placeholder="Enter your email / nickname" 
                            required
                        >
                    </div>

                    <div class="form-group">
                        <label for="password">Password</label>
                        <input 
                            type="password" 
                            id="password" 
                            class="form-control" 
                            placeholder="Enter your password" 
                            required
                        >
                    </div>

                    <div class="form-actions">
                        <button type="submit" class="btn btn-primary">Sign In</button>
                    </div>

                    <div class="form-footer">
                        <p>
                            Don't have an account? 
                            <a href="/sign-up" data-link>Sign Up</a>
                        </p>
                            <a href="/" data-link>Go Home</a>
                    </div>

                    <div id="showError" class="error-message"></div>
                </form>
            </div>
        `;
  }

  async init() {
    const signInForm = document.getElementById("sign-in-form");

    signInForm.addEventListener("submit", async (e) => {
      e.preventDefault(); // Prevent default form submission

      const login = document.getElementById("login")?.value;
      const password = document.getElementById("password")?.value;

      // Validate inputs
      if (!login || !password) {
        document.getElementById("showError").innerHTML =
          "Please enter both email and password";
        return;
      }

      await signIn(login, password);
    });

    const loginInput = document.getElementById("login");
    const passwordInput = document.getElementById("password");
    let showErr = document.getElementById("showError");

    loginInput.addEventListener("input", () => {
      // Basic email validation
      if (loginInput.value?.length < 3) {
        showErr.innerHTML = "Invalid Email / Nickname";
      } else {
        showErr.innerHTML = "";
      }
    });

    passwordInput.addEventListener("input", () => {
      // Basic password validation
      if (passwordInput.value.length < 8) {
        showErr.innerHTML = "Password must be at least 6 characters long";
      } else {
        showErr.innerHTML = "";
      }
    });
  }
}
