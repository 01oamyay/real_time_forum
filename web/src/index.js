import Home from "./views/HomeView.js";
import SignIn from "./views/SignInView.js";
import SignUp from "./views/SignUpView.js";
import CreatePost from "./views/CreatePostView.js";
import Post from "./views/PostView.js";
import NavBar from "./views/NavBarView.js";
import SideBar from "./views/SideBarView.js";
import Utils from "./pkg/Utils.js";
import fetcher from "./pkg/fetcher.js";
import UsersListView from "./views/UsersListView.js";
import WS from "./pkg/WS.js";
import ChatView from "./views/ِChatView.js";

let chatInstance = null;

const pathToRegex = (path) =>
  new RegExp("^" + path.replace(/\//g, "\\/").replace(/:\w+/g, "(.+)") + "$");

const roles = {
  guest: 0,
  user: 1,
};

let ws = new WS();

const getParams = (match) => {
  const values = match.result.slice(1);
  const keys = Array.from(match.route.path.matchAll(/:(\w+)/g)).map(
    (result) => result[1]
  );

  return Object.fromEntries(
    keys.map((key, i) => {
      return [key, values[i]];
    })
  );
};

const navigateTo = (url) => {
  history.pushState(null, null, url);
  router();
};

const router = async () => {
  const routes = [
    { path: "/", view: Home, minRole: roles.user, style: "main-content" },
    { path: "/sign-in", view: SignIn, minRole: roles.guest, style: "auth" },
    { path: "/sign-up", view: SignUp, minRole: roles.guest, style: "auth" },
    {
      path: "/create-post",
      view: CreatePost,
      minRole: roles.user,
      style: "create-post",
    },
    { path: "/post/:postID", view: Post, minRole: roles.user, style: "post" },
    {
      path: "/chat/:userID",
      view: ChatView,
      minRole: roles.user,
      style: "chat",
    },
  ];

  const potentialMatches = routes.map((route) => {
    return {
      route: route,
      result: location.pathname.match(pathToRegex(route.path)),
    };
  });

  const checker = await fetcher.checkToken();
  if (checker && !checker.checker) {
    localStorage.setItem("role", roles.guest);
    localStorage.removeItem("id");
  }

  const user = Utils.getUser();

  if (!user.role) {
    user.role = roles.guest;
    localStorage.setItem("role", user.role);
  }

  let match = potentialMatches.find(
    (potentialMatches) => potentialMatches.result !== null
  );

  if (!match) {
    Utils.showError(404, "The page you requested does not exist");
    return;
  }

  const isLogged = await fetcher.isLoggedIn();
  if (
    isLogged &&
    (match.route.path == "/sign-in" || match.route.path == "/sign-up")
  ) {
    navigateTo("/");
    return;
  }

  const view = new match.route.view(getParams(match), user);

  if (match.route.view == ChatView) {
    if (chatInstance) {
      chatInstance.destroy();
    }

    chatInstance = view;
  }

  // Remove previous view-specific styles
  view.removeStyles();

  // Add new view-specific style
  if (match.route.style) {
    view.addStyle(match.route.style);
  }

  if (user.role < match.route.minRole) {
    navigateTo("/sign-in");
    return;
  }

  if (
    match.route.view === Home ||
    match.route.view === Post ||
    match.route.view === CreatePost ||
    match.route.view === ChatView
  ) {
    // Load Navbar
    const NavBarView = new NavBar(null, user);
    document.querySelector("#navbar").innerHTML = await NavBarView.getHtml();
    NavBarView.init();

    view.addStyle("navbar");
    view.addStyle("main-content");
    view.addStyle("post-card");

    let sideBarHTML = "";
    let userListHTML = "";
    // Load Sidebar
    let SideBarView;
    // Load Userlist
    let UserListView;
    view.addStyle("user_sidebar");
    if (match.route.view === Home || match.route.view == ChatView) {
      view.addStyle("sidebar");

      ws?.init();

      if (match.route.view == Home) {
        SideBarView = new SideBar(null, user);
        sideBarHTML = await SideBarView.getHtml();
      }
    }
    UserListView = new UsersListView(null, user);
    userListHTML = await UserListView.getHtml();
    document.querySelector("#app").innerHTML = `<div id="toast">
      <div id="toast-message"></div>`;
    document.querySelector("#app").innerHTML += sideBarHTML;
    document.querySelector("#app").innerHTML += await view.getHtml();
    document.querySelector("#app").innerHTML += userListHTML;

    SideBarView?.init();
    UserListView?.init();
  } else {
    // Clear navbar and sidebar if not HomeView
    document.querySelector("#navbar").innerHTML = "";
    document.querySelector("#app").innerHTML = await view.getHtml();
  }

  view.init();
};

window.addEventListener("popstate", router);

window.addEventListener("storage", () => {
  const user = Utils.getUser();
  if (user.id == null) {
    location.reload();
  }
});

document.addEventListener("DOMContentLoaded", () => {
  document.body.addEventListener("click", (e) => {
    if (e.target.matches("[data-link]")) {
      e.preventDefault();
      navigateTo(e.target.href);
    } else if (e.target.closest("[data-link]")) {
      e.preventDefault();
      navigateTo(e.target.closest("[data-link]")?.href);
    }
  });

  document.addEventListener("ws-closing", () => {
    Utils.showToast(
      "Websocket is either closed or in closing state, the page will be refreshed",
      "error"
    );

    // refresh the page after 3 seconds
    setTimeout(() => {
      location.reload();
    }, 3000);
  });

  router();
});

export default { navigateTo, roles };
