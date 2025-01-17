import AbstractView from "./AbstractView.js";
import redirect from "../index.js";
import fetcher from "../pkg/fetcher.js";

const signup = async (email, nickname, password, rePassword) => {
  let body = {
    email: email,
    nickname: nickname,
    password: password,
    cfmpsw: rePassword,
  };

  const data = await fetcher.post("/api/signup", body);
  if (data && data.msg !== undefined) {
    let showErr = document.getElementById("showError");
    showErr.innerHTML = data.msg;
    return;
  }
  redirect.navigateTo("/sign-in");
};

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Sign Up");
  }

  async getHtml() {
    return `
            <div class="auth-wrapper">
                <form id="sign-up-form" class="auth-form">
                    <h2 class="form-title">Create Your Account</h2>
                    
                    <div class="form-group">
                        <label for="email">Email Address</label>
                        <input 
                            type="email" 
                            id="email" 
                            class="form-control" 
                            placeholder="Enter your email" 
                            required
                        >
                        <div class="validation-message" id="email-error"></div>
                    </div>

                    <div class="form-group">
                        <label for="nickname">Nickname</label>
                        <input 
                            type="text" 
                            id="nickname" 
                            class="form-control" 
                            placeholder="Choose a nickname" 
                            required
                            minlength="3"
                            maxlength="20"
                        >
                        <div class="validation-message" id="nickname-error"></div>
                    </div>

                    <div class="form-group">
                        <label for="password">Password</label>
                        <input 
                            type="password" 
                            id="password" 
                            class="form-control" 
                            placeholder="Create a password" 
                            required
                            minlength="6"
                        >
                        <div class="validation-message" id="password-error"></div>
                    </div>

                    <div class="form-group">
                        <label for="rePassword">Confirm Password</label>
                        <input 
                            type="password" 
                            id="rePassword" 
                            class="form-control" 
                            placeholder="Repeat your password" 
                            required
                        >
                        <div class="validation-message" id="rePassword-error"></div>
                    </div>

                    <div class="form-actions">
                        <button type="submit" class="btn btn-primary">Sign Up</button>
                    </div>

                    <div class="form-footer">
                        <p>
                            Already have an account? 
                            <a href="/sign-in" data-link>Sign In</a>
                        </p>
                          <a href="/" data-link>Go Home</a>
                    </div>

                    <div id="showError" class="error-message"></div>
                </form>
            </div>
        `;
  }

  async init() {
    const signUpForm = document.getElementById("sign-up-form");

    // Input elements
    const emailInput = document.getElementById("email");
    const nicknameInput = document.getElementById("nickname");
    const passwordInput = document.getElementById("password");
    const rePasswordInput = document.getElementById("rePassword");

    // Error message elements
    const emailError = document.getElementById("email-error");
    const nicknameError = document.getElementById("nickname-error");
    const passwordError = document.getElementById("password-error");
    const rePasswordError = document.getElementById("rePassword-error");

    // Email validation
    const validateEmail = (email) => {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!email) {
        emailError.textContent = "Email is required";
        return false;
      }
      if (!emailRegex.test(email)) {
        emailError.textContent = "Please enter a valid email address";
        return false;
      }
      emailError.textContent = "";
      return true;
    };

    // nickname validation
    const validatenickname = (nickname) => {
      if (!nickname) {
        nicknameError.textContent = "nickname is required";
        return false;
      }

      if (nickname.length < 3) {
        nicknameError.textContent =
          "nickname must be at least 3 characters long";
        return false;
      }
      if (nickname.length > 20) {
        nicknameError.textContent = "nickname must be less than 20 characters";
        return false;
      }
      nicknameError.textContent = "";
      return true;
    };

    // Password validation
    const validatePassword = (password) => {
      if (!password) {
        passwordError.textContent = "Password is required";
        return false;
      }
      if (password.length < 8) {
        passwordError.textContent =
          "Password must be at least 6 characters long";
        return false;
      }
      passwordError.textContent = "";
      return true;
    };

    // Confirm password validation
    const validateConfirmPassword = (password, confirmPassword) => {
      if (!confirmPassword) {
        rePasswordError.textContent = "Please confirm your password";
        return false;
      }
      if (password !== confirmPassword) {
        rePasswordError.textContent = "Passwords do not match";
        return false;
      }
      rePasswordError.textContent = "";
      return true;
    };

    // Real-time validation
    emailInput.addEventListener("input", () => validateEmail(emailInput.value));
    nicknameInput.addEventListener("input", () =>
      validatenickname(nicknameInput.value)
    );
    passwordInput.addEventListener("input", () =>
      validatePassword(passwordInput.value)
    );
    rePasswordInput.addEventListener("input", () =>
      validateConfirmPassword(passwordInput.value, rePasswordInput.value)
    );

    // Form submission
    signUpForm.addEventListener("submit", async (e) => {
      e.preventDefault();

      // Validate all fields
      const isEmailValid = validateEmail(emailInput.value);
      const isnicknameValid = validatenickname(nicknameInput.value);
      const isPasswordValid = validatePassword(passwordInput.value);
      const isConfirmPasswordValid = validateConfirmPassword(
        passwordInput.value,
        rePasswordInput.value
      );

      // Only proceed if all validations pass
      if (
        isEmailValid &&
        isnicknameValid &&
        isPasswordValid &&
        isConfirmPasswordValid
      ) {
        await signup(
          emailInput.value,
          nicknameInput.value,
          passwordInput.value,
          rePasswordInput.value
        );
      }
    });
  }
}
