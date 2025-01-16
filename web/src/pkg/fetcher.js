import Utils from "./Utils.js";
import redirect from "../index.js";

const fetcher = {
  get: async (path, body) => {
    return makeRequest(path, body, "GET");
  },
  post: async (path, body) => {
    return makeRequest(path, body, "POST");
  },
  checkToken: async () => {
    //http://${API_HOST_NAME}
    const url = `/api/is-valid`;
    const options = {
      mode: "cors",
      method: "GET",
    };

    const response = await fetch(url, options).catch((e) => {
      console.log(e);
      Utils.showError(503);
      return;
    });
    var responseBody;
    try {
      responseBody = await response.json();
    } catch {
      console.log("some unexpected error: with json");
      return;
    }
    if (!response.ok) {
      Utils.showError(response.status, responseBody.msg);
      return responseBody;
    }

    return responseBody;
  },
  isLoggedIn: async () => {
    const url = `/api/is-valid`;
    const options = {
      mode: "cors",
      method: "GET",
    };

    const response = await fetch(url, options);
    if (!response.ok) {
      return false;
    } else {
      const responseBody = await response.json();
      return responseBody.checker;
    }
  },
};

const makeRequest = async (path, body, method) => {
  //http://${API_HOST_NAME}
  const url = `${path}`;
  const options = {
    mode: "cors",
    method: method,
    body: JSON.stringify(body),
  };

  const response = await fetch(url, options).catch((e) => {
    console.log(e);
    Utils.showError(503);
    return;
  });

  var responseBody;
  try {
    responseBody = await response.json();
  } catch {
    return;
  }

  if (response.status == 401 || response.status == 403) {
    Utils.logOut();
    redirect.navigateTo("/sign-in");
    return responseBody;
  }
  if (response.status == 404) {
    return { status: 404, msg: "" };
  }
  if (response.status == 400) {
    return responseBody;
  }
  if (!response.ok) {
    Utils.showError(response.status, responseBody.msg);
    return responseBody;
  }
  return responseBody;
};

export default fetcher;
