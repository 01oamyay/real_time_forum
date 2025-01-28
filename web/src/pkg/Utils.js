const showError = (status, message) => {
  const app = document.querySelector("#app");

  const titles = {
    400: "400 Bad Request",
    401: "401 Unauthorized",
    403: "403 Forbidden",
    404: "404 Not Found",
    405: "405 Method Not Allowed",
    429: "429 Too Many Requests",
    500: "500 Internal Server Error",
    503: "503 Service Unavailable",
  };

  app.innerHTML = `
        <div class="errorDiv">
        <h1>${titles[status]}</h1><br>
        <h2>${message || ""}</h2>
        </div>
    `;
};

const getUser = () => {
  return {
    id: localStorage.getItem("id"),
    role: localStorage.getItem("role"),
  };
};

const logOut = () => {
  localStorage.removeItem("id");
  localStorage.removeItem("role");
 
  fetch(`/api/signout`, {
    mode: "cors",
    method: "POST",
  }).then(() => window.location.reload());
};

const parseJwt = (token) => {
  var base64Url = token.split(".")[1];
  var base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  var jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );

  return JSON.parse(jsonPayload);
};

function showToast(msg, type = "error") {
  const toast = document.getElementById("toast");
  toast.innerHTML = msg;
  toast.classList.add(type);
  toast.classList.add("show");

  setTimeout(() => {
    toast.classList.remove("show");
    toast.classList.remove(type);
  }, 5000);
}

export default { showError, getUser, logOut, parseJwt, showToast };
