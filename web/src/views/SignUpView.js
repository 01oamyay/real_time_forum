import AbstractView from "./AbstractView.js";
import redirect from "../index.js";
import fetcher from "../pkg/fetcher.js";

const signup = async (user) => {

  const data = await fetcher.post("/api/signup", user);
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
                        <label for="firstName">First Name</label>
                        <input 
                            type="text" 
                            id="firstName" 
                            class="form-control" 
                            placeholder="Choose a first Name" 
                            required
                            minlength="3"
                            maxlength="20"
                        >
                        <div class="validation-message" id="firstName-error"></div>
                    </div>        
                    <div class="form-group">
                        <label for="lastName">Last Name</label>
                        <input 
                            type="text" 
                            id="lastName" 
                            class="form-control" 
                            placeholder="Choose a last Name" 
                            required
                            minlength="3"
                            maxlength="20"
                        >
                        <div class="validation-message" id="lastName-error"></div>
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
                        <label for="age">Age</label>
                        <input 
                            type="number" 
                            id="age"
                            class="form-control" 
                            placeholder="Choose a age" 
                            required
                        >
                        <div class="validation-message" id="age-error"></div>
                    </div>
                    <div class="form-group">
                        <label for="gender">Gender</label>
                        <select 
                            id="gender" 
                            name="gender"
                            class="form-control" 
                            placeholder="Choose a gender"
                            required
                        >
                          <option value="male">Male</option>
                          <option value="female">Female</option>
                        </select>
                        <div class="validation-message" id="gender-error"></div>
                    </div>
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
    const firstNameInput = document.getElementById("firstName");
    const lastNameInput = document.getElementById("lastName");
    const nicknameInput = document.getElementById("nickname");
    const ageInput = document.getElementById("age");
    const genderInput = document.getElementById("gender");
    const emailInput = document.getElementById("email");
    const passwordInput = document.getElementById("password");
    const rePasswordInput = document.getElementById("rePassword");

    // Error message elements
    const firstNameError = document.getElementById("firstName-error");
    const lastNameError = document.getElementById("lastName-error");
    const nicknameError = document.getElementById("nickname-error");
    const ageError = document.getElementById("age-error");
    const genderError = document.getElementById("gender-error");
    const emailError = document.getElementById("email-error");
    const passwordError = document.getElementById("password-error");
    const rePasswordError = document.getElementById("rePassword-error");

    // First name validation
    const validateFirstName = (firstName) => {
      if (!firstName) {
        firstNameError.textContent = "First name is required";
        return false;
      }
      if (firstName.length < 3) {
        firstNameError.textContent =
          "First name must be at least 3 characters long";
        return false;
      }
      if (firstName.length > 20) {
        firstNameError.textContent =
          "First name must be less than 20 characters";
        return false;
      }
      firstNameError.textContent = "";
      return true;
    };

    // Last name validation
    const validateLastName = (lastName) => {
      if (!lastName) {
        lastNameError.textContent = "Last name is required";
        return false;
      }
      if (lastName.length < 3) {
        lastNameError.textContent =
          "Last name must be at least 3 characters long";
        return false;
      }
      if (lastName.length > 20) {
        lastNameError.textContent = "Last name must be less than 20 characters";
        return false;
      }
      lastNameError.textContent = "";
      return true;
    };

    // Nickname validation
    const validateNickname = (nickname) => {
      if (!nickname) {
        nicknameError.textContent = "Nickname is required";
        return false;
      }
      if (nickname.length < 3) {
        nicknameError.textContent =
          "Nickname must be at least 3 characters long";
        return false;
      }
      if (nickname.length > 20) {
        nicknameError.textContent = "Nickname must be less than 20 characters";
        return false;
      }
      nicknameError.textContent = "";
      return true;
    };

    // Age validation
    const validateAge = (age) => {
      if (!age) {
        ageError.textContent = "Age is required";
        return false;
      }
      if (age < 18) {
        ageError.textContent = "Age must be at least 18";
        return false;
      }
      ageError.textContent = "";
      return true;
    };

    // Gender validation
    const validateGender = (gender) => {
      if (!gender) {
        genderError.textContent = "Gender is required";
        return false;
      }
      genderError.textContent = "";
      return true;
    };

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

    // Password validation
    const validatePassword = (password) => {
      if (!password) {
        passwordError.textContent = "Password is required";
        return false;
      }
      if (password.length < 6) {
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
    firstNameInput.addEventListener("input", () =>
      validateFirstName(firstNameInput.value)
    );
    lastNameInput.addEventListener("input", () =>
      validateLastName(lastNameInput.value)
    );
    nicknameInput.addEventListener("input", () =>
      validateNickname(nicknameInput.value)
    );
    ageInput.addEventListener("input", () => validateAge(ageInput.value));
    genderInput.addEventListener("input", () =>
      validateGender(genderInput.value)
    );
    emailInput.addEventListener("input", () => validateEmail(emailInput.value));
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
      const isFirstNameValid = validateFirstName(firstNameInput.value);
      const isLastNameValid = validateLastName(lastNameInput.value);
      const isNicknameValid = validateNickname(nicknameInput.value);
      const isAgeValid = validateAge(ageInput.value);
      const isGenderValid = validateGender(genderInput.value);
      const isEmailValid = validateEmail(emailInput.value);
      const isPasswordValid = validatePassword(passwordInput.value);
      const isConfirmPasswordValid = validateConfirmPassword(
        passwordInput.value,
        rePasswordInput.value
      );

      // Only proceed if all validations pass
      if (
        isFirstNameValid &&
        isLastNameValid &&
        isNicknameValid &&
        isAgeValid &&
        isGenderValid &&
        isEmailValid &&
        isPasswordValid &&
        isConfirmPasswordValid
      ) {
        await signup({
          email: emailInput.value,
          nickname: nicknameInput.value,
          firstName: firstNameInput.value,
          lastName: lastNameInput.value,
          gender: genderInput.value,
          age: Number(ageInput.value),
          password: passwordInput.value,
          cfmpsw: rePasswordInput.value,
        });
      }
    });
  }
}
