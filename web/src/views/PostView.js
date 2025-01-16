import fetcher from "../pkg/fetcher.js";
import Utils from "../pkg/Utils.js";
import AbstractView from "./AbstractView.js";

const getPostPath = "/api/post/";
const sendCommentPath = "/api/comment/create";

let likeListenerSet = false; // Flag to track like button listener
let dislikeListenerSet = false; // Flag to track dislike button listener

const getPost = async (postID) => {
  const post = await fetcher.get(getPostPath + postID);
  if (post && post.msg != undefined) {
    Utils.showError(post.status, post.msg);
    return;
  }
  if (post) {
    document.getElementById("post-title").innerText = post.title;
    const userEl = document.getElementById("post-user-id");
    userEl.innerText = "Author: " + post.username;

    const formattedCategories = post?.categories?.map((cat) => " #" + cat);
    document.getElementById("post-tags").innerText =
      "Categories:" + formattedCategories;

    document.getElementById("post-data").innerText = post.data;
    document.getElementById("post-like-inner").innerText = post.likes;
    document.getElementById("post-dislike-inner").innerText = post.dislikes;

    const likeBtn = document.getElementById("post-like");
    const dislikeBtn = document.getElementById("post-dislike");

    // Define named functions for event listeners
    const handleLike = () => {
      votePost(postID, 1);
    };

    const handleDislike = () => {
      votePost(postID, 0);
    };

    // Set up like button listener if not already set
    if (!likeListenerSet) {
      likeBtn.addEventListener("click", handleLike);
      likeListenerSet = true; // Mark as set
    }

    // Set up dislike button listener if not already set
    if (!dislikeListenerSet) {
      dislikeBtn.addEventListener("click", handleDislike);
      dislikeListenerSet = true; // Mark as set
    }

    const commentsDoc = document.getElementById("comments");
    commentsDoc.innerHTML = ""; // Clear previous comments

    if (post.comments.length > 0) {
      const commentText = document.createElement("h3");
      commentText.innerText = "Comments: ";
      commentsDoc.append(commentText);
    }

    for (let i = post.comments.length - 1; i >= 0; i--) {
      const comment = post.comments[i];
      const el = drawComments(comment);
      commentsDoc.append(el);
    }
  }
};

const votePost = async (postID, likeType) => {
  const path = "/api/post/vote";
  const body = {
    post_id: parseInt(postID),
    vote: likeType,
  };
  const data = await fetcher.post(path, body);
  if (data && data.msg) {
    return;
  }
  await getPost(postID);
};

const voteComment = async (commentID, likeType) => {
  const path = "/api/comment/vote";
  const body = {
    comment_id: parseInt(commentID),
    vote: likeType,
  };
  const data = await fetcher.post(path, body);
  if (data && data.msg) {
    return;
  }

  const postIDFromURL = window.location.pathname.split("/").pop();

  if (postIDFromURL) {
    await getPost(postIDFromURL);
  } else {
    console.error("Could not find post ID");
  }
};

const drawComments = (comment) => {
  const el = document.createElement("div");
  el.classList.add("card");

  const authorEl = document.createElement("p");
  authorEl.classList.add("card-header");
  authorEl.innerText = "Author: " + comment.username;

  const body = document.createElement("div");
  body.classList.add("card-body");

  const dataEl = document.createElement("p");
  dataEl.classList.add("card-text");
  dataEl.innerText = comment.data;
  body.append(dataEl);

  // Like Icon (Thumbs Up)
  const likeIcon = document.createElementNS(
    "http://www.w3.org/2000/svg",
    "svg"
  );
  likeIcon.setAttribute("xmlns", "http://www.w3.org/2000/svg");
  likeIcon.setAttribute("viewBox", "0 0 24 24");
  likeIcon.setAttribute("fill", "none");
  likeIcon.setAttribute("stroke", "currentColor");
  likeIcon.setAttribute("stroke-width", "2");
  likeIcon.classList.add("comment-stats-icon");
  likeIcon.innerHTML = `
      <path
        d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    `;

  // Dislike Icon (Thumbs Down)
  const dislikeIcon = document.createElementNS(
    "http://www.w3.org/2000/svg",
    "svg"
  );
  dislikeIcon.setAttribute("xmlns", "http://www.w3.org/2000/svg");
  dislikeIcon.setAttribute("viewBox", "0 0 24 24");
  dislikeIcon.setAttribute("fill", "none");
  dislikeIcon.setAttribute("stroke", "currentColor");
  dislikeIcon.setAttribute("stroke-width", "2");
  dislikeIcon.classList.add("comment-stats-icon");
  dislikeIcon.innerHTML = `
  <path 
    d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zM17 2h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2h-3"
    stroke-linecap="round"
    stroke-linejoin="round"
  />`;

  let likeButton = document.createElement("button");
  likeButton.className = "btn comment-like";
  likeButton.id = "comment-like";
  likeIcon.innerText = comment.likes;

  let likeCount = document.createElement("p");
  likeCount.className = "comment-like-count";
  likeCount.innerText = comment.likes;

  likeButton.appendChild(likeIcon);
  likeButton.appendChild(likeCount);

  let dislikeButton = document.createElement("button");
  dislikeButton.className = "btn comment-dislike";
  dislikeButton.id = "comment-dislike";
  dislikeIcon.innerText = comment.dislikes;

  let dislikeCount = document.createElement("p");
  dislikeCount.className = "comment-dislike-count";
  dislikeCount.innerText = comment.dislikes;

  dislikeButton.appendChild(dislikeIcon);
  dislikeButton.appendChild(dislikeCount);

  if (comment?.vote_status == 1) {
    if (likeButton) {
      likeButton.classList.add("active");
      dislikeButton.classList.remove("active");
    }
  }

  if (comment?.vote_status == 2) {
    if (dislikeButton) {
      dislikeButton.classList.add("active");
      likeButton.classList.remove("active");
    }
  }

  const votes = document.createElement("div");
  votes.classList.add("comment-votes");
  votes.appendChild(likeButton);
  votes.appendChild(dislikeButton);

  likeButton.addEventListener("click", () => {
    voteComment(comment.comment_id, 1);
  });
  dislikeButton.addEventListener("click", () => {
    voteComment(comment.comment_id, 0);
  });

  el.append(authorEl);
  el.append(body);
  el.append(votes);
  return el;
};

const sendComment = async (comment, postID) => {
  let body = {
    data: comment,
    post_id: parseInt(postID),
  };
  const data = await fetcher.post(sendCommentPath, body);
  if (data && data.msg !== undefined) {
    let showErr = document.getElementById("showError");
    showErr.innerHTML = data.msg;
    return;
  }
  await getPost(postID);
};

export default class extends AbstractView {
  constructor(params, user) {
    super(params);
    this.user = user;
    this.setTitle("Post");
  }

  async getHtml() {
    const isAuthorized = Boolean(this.user.id);
    return `
    
            <div class="post-container">
                <div class="post-details">
                    <h3 id="post-title"></h3>
                    <div id="post-user">
                        <span id="post-user-id"></span>
                    </div>
                    <h5 id="post-tags"></h5>
                    <div class="post-card-content">
                        <p id="post-data"></p>
                    </div>
                    <div class="post-actions">
                    <button class="btn post-like" id="post-like">
                    <svg xlmns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="post-stats-icon">
                      <path
                        d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3zM7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        />
                    </svg>
                    <p id="post-like-inner"></p>
                    </button>
                    <button class="btn post-dislike" id="post-dislike">
                      <svg xlmns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="post-stats-icon">
                          <path 
                            d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zM17 2h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2h-3"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                          />
                      </svg>
                    <p id="post-dislike-inner"></p>
                    </button>
                    </div>
                </div>
            </div>

            ${
              isAuthorized
                ? `
            <div class="comment-section">
                <form id="comment-form" class="comment-form">
                    <h3>Leave a comment here:</h3>
                    <textarea 
                        id="comment-input" 
                        class="form-control" 
                        rows="3" 
                        placeholder="Leave a comment"
                    ></textarea>
                    <button type="submit" class="btn btn-primary">Send</button>
                    <div class="error" id="showError"></div>
                </form>
            </div>
            `
                : ""
            }

            <div id="comments" class="comments-container">
                <!-- Comments will be dynamically populated here -->
            </div>
        `;
  }

  async init() {
    const postID = this.params.postID;
    await getPost(postID);

    // Comment submission for authorized users
    const commentForm = document.getElementById("comment-form");
    commentForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const commentInput = document.getElementById("comment-input");
      const comment = commentInput.value.trim();
      if (comment) {
        await sendComment(comment, postID);
        commentInput.value = ""; // Clear the input after submission
      }
    });
  }
}
