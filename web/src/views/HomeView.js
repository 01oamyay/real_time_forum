import AbstractView from "./AbstractView.js";
import fetcher from "../pkg/fetcher.js";
import Utils from "../pkg/Utils.js";

const path = `/api/posts/`;
let offset = 0;
let limit = 20;
let ended = false;

const getPostsByCategory = async (category) => {
  let posts = await fetcher.get(
    `${path}${category}?limit=${limit}&offset=${offset}`
  );
  const postsDoc = document.querySelector(".posts-grid");

  if (posts && posts?.status == "400" && offset > 0) {
    ended = true;
    return;
  }

  if (posts && posts.msg !== undefined) {
    postsDoc.innerHTML = "<p>No posts</p>";
    return;
  }
  if (posts?.length) {
    posts.sort((a, b) => a.post_id - b.post_id);
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
    offset += posts.length;
  } else {
    if (offset > 0) {
      ended = true;
      return;
    }
    postsDoc.innerHTML = "<p>No posts</p>"; // Display "No posts" message
  }
};

const newPostElement = (post) => {
  const postCard = document.createElement("div");
  postCard.classList.add("post-card");
  postCard.setAttribute("data-post-id", post.post_id);
  // Post Card Header
  const postCardHeader = document.createElement("div");
  postCardHeader.classList.add("post-card-header");
  const headerContent = document.createElement("div");
  const authorName = document.createElement("div");
  authorName.classList.add("post-author-name");
  authorName.textContent = post.username;

  headerContent.appendChild(authorName);
  postCardHeader.appendChild(headerContent);

  // Post Content
  const postContent = document.createElement("div");
  postContent.classList.add("post-content");

  const postTitle = document.createElement("h3");
  postTitle.classList.add("post-title");
  postTitle.textContent = post.title;

  const postExcerpt = document.createElement("p");
  postExcerpt.classList.add("post-excerpt");
  postExcerpt.textContent = post.data.substring(0, 200) + "...";

  const postCategory = document.createElement("div");
  postCategory.classList.add("post-category");
  postCategory.style.color = "#3182ce";
  postCategory.textContent = post.categories
    ?.filter((c) => c != "General")
    ?.join(", ");

  postContent.appendChild(postTitle);
  postContent.appendChild(postExcerpt);
  postContent.appendChild(postCategory);

  // Post Actions
  const postActions = document.createElement("div");
  postActions.classList.add("post-card-actions");

  const postStats = document.createElement("div");
  postStats.classList.add("post-stats");

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
  likeIcon.classList.add("post-stats-icon");
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
  dislikeIcon.classList.add("post-stats-icon");
  dislikeIcon.innerHTML = `
<path 
  d="M10 15v4a3 3 0 0 0 3 3l4-9V2H5.72a2 2 0 0 0-2 1.7l-1.38 9a2 2 0 0 0 2 2.3zM17 2h3a2 2 0 0 1 2 2v7a2 2 0 0 1-2 2h-3"
  stroke-linecap="round"
  stroke-linejoin="round"
/>`;

  // Comment Icon (Speech Bubble with Text)
  const commentIcon = document.createElementNS(
    "http://www.w3.org/2000/svg",
    "svg"
  );
  commentIcon.setAttribute("xmlns", "http://www.w3.org/2000/svg");
  commentIcon.setAttribute("viewBox", "0 0 24 24");
  commentIcon.setAttribute("fill", "none");
  commentIcon.setAttribute("stroke", "currentColor");
  commentIcon.setAttribute("stroke-width", "2");
  commentIcon.classList.add("post-stats-icon");
  commentIcon.innerHTML = `
  <path 
    d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
    stroke-linejoin="round"
  />
  <line x1="8" y1="9" x2="16" y2="9" stroke-linecap="round" />
  <line x1="8" y1="13" x2="12" y2="13" stroke-linecap="round" />
`;

  // Like count
  const likeCount = document.createElement("span");
  likeCount.textContent = post.likes || 0;
  likeCount.classList.add("post-stats-count");

  // Dislike count
  const dislikeCount = document.createElement("span");
  dislikeCount.textContent = post.dislikes || 0;
  dislikeCount.classList.add("post-stats-count");

  // Comment count
  const commentCount = document.createElement("span");
  commentCount.textContent = post.comments_count || 0;
  commentCount.classList.add("post-stats-count");

  // Like Button
  const likeButton = document.createElement("button");
  likeButton.classList.add("post-stats-button", "like-button");
  likeButton.appendChild(likeIcon);
  likeButton.appendChild(likeCount);

  likeButton.addEventListener("click", () => {
    votePost(post.post_id, 1);
  });

  // Dislike Button
  const dislikeButton = document.createElement("button");
  dislikeButton.classList.add("post-stats-button", "dislike-button");
  dislikeButton.appendChild(dislikeIcon);
  dislikeButton.appendChild(dislikeCount);

  dislikeButton.addEventListener("click", () => {
    votePost(post.post_id, 0);
  });

  const commentButton = document.createElement("button");
  commentButton.classList.add("post-stats-button");
  commentButton.appendChild(commentIcon);
  commentButton.appendChild(commentCount);

  commentButton.addEventListener("click", () => {
    window.location.href = `/post/${post.post_id}`;
  });

  // check vote status
  if (post?.vote_status == 1) {
    if (likeButton) {
      likeButton.classList.add("active");
      dislikeButton.classList.remove("active");
    }
  }
  if (post?.vote_status == 2) {
    if (dislikeButton) {
      dislikeButton.classList.add("active");
      likeButton.classList.remove("active");
    }
  }

  // Append buttons to stats
  postStats.appendChild(likeButton);
  postStats.appendChild(dislikeButton);
  postStats.appendChild(commentButton);

  postActions.appendChild(postStats);

  // Assemble post card
  postCard.appendChild(postCardHeader);
  postCard.appendChild(postContent);
  postCard.appendChild(postActions);

  return postCard;
};

const votePost = async (postId, voteType) => {
  const path = "/api/post/vote";
  const body = {
    post_id: parseInt(postId),
    vote: voteType,
  };

  try {
    const data = await fetcher.post(path, body);
    if (data && data.msg) {
      console.error(data.msg);
      return;
    }

    // Fetch the updated post details
    await updatePostInList(postId);
  } catch (error) {
    console.error("Error voting on post:", error);
  }
};

const updatePostInList = async (postId) => {
  try {
    const updatedPost = await fetcher.get(`/api/post/${postId}`);
    const postCard = document.querySelector(
      `.post-card[data-post-id="${postId}"]`
    );

    if (!postCard || !updatedPost) return;

    // Most robust method
    const updateCount = (buttonClass, count) => {
      const button = postCard.querySelector(`.${buttonClass}`);
      if (button) {
        const countEl = button.querySelector(".post-stats-count");
        if (countEl) {
          countEl.textContent = count || 0;
        }
      }
    };

    updateCount("like-button", updatedPost.likes);
    updateCount("dislike-button", updatedPost.dislikes);

    const likeButton = postCard.querySelector(".like-button");
    const dislikeButton = postCard.querySelector(".dislike-button");

    if (updatedPost?.vote_status == 1) {
      likeButton?.classList?.add("active");
      dislikeButton?.classList?.remove("active");
    } else if (updatedPost?.vote_status == 2) {
      dislikeButton?.classList?.add("active");
      likeButton?.classList?.remove("active");
    } else {
      dislikeButton?.classList?.remove("active");
      likeButton?.classList?.remove("active");
    }
  } catch (error) {
    console.error("Error updating post:", error);
  }
};

const getUserPosts = async (userID) => {
  const path = `/api/profile/posts/${userID}?limit=${limit}&offset=${offset}`;
  const posts = await fetcher.get(path);

  if (posts && posts?.status == "404" && offset > 0) {
    ended = true;
    return;
  }

  const postsDoc = document.querySelector(".posts-grid");

  if (posts && posts.msg !== undefined) {
    postsDoc.innerHTML = "<p>No posts</p>";
    return;
  }
  if (posts?.length) {
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
    offset += 10;
  } else {
    if (offset > 0) {
      ended = true;
      return;
    }
    postsDoc.innerHTML = "<p>No posts</p>"; // Display "No posts" message
  }
};

const getUserLikedPosts = async (userID) => {
  const path = `/api/profile/liked-posts/${userID}?limit=${limit}&offset=${offset}`;
  const posts = await fetcher.get(path);

  if (posts && posts?.status == "404" && offset > 0) {
    ended = true;
    return;
  }

  const postsDoc = document.querySelector(".posts-grid");

  if (posts && posts.msg !== undefined) {
    postsDoc.innerHTML = "<p>No posts</p>";
    return;
  }
  if (posts?.length) {
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
    offset += 10;
  } else {
    if (offset > 0) {
      ended = true;
      return;
    }
    postsDoc.innerHTML = "<p>No posts</p>"; // Display "No posts" message
  }
};

const getUserDislikedPosts = async (userID) => {
  const path = `/api/profile/disliked-posts/${userID}?limit=${limit}&offset=${offset}`;
  const posts = await fetcher.get(path);

  if (posts && posts?.status == "404" && offset > 0) {
    ended = true;
    return;
  }

  const postsDoc = document.querySelector(".posts-grid");

  if (posts && posts.msg !== undefined) {
    postsDoc.innerHTML = "<p>No posts</p>";
    return;
  }
  if (posts?.length) {
    for (let i = posts.length - 1; i >= 0; i--) {
      const post = posts[i];
      const el = newPostElement(post);
      postsDoc.append(el);
    }
    offset += 10;
  } else {
    if (offset > 0) {
      ended = true;
      return;
    }
    postsDoc.innerHTML = "<p>No posts</p>"; // Display "No posts" message
  }
};

let isThrottled = false;
let category = "General";

async function handleScroll() {
  if (ended) {
    const postsContainer = document.getElementById("posts");
    postsContainer?.removeEventListener("scroll", handleScroll); // Remove the event listener
    return;
  }

  if (isThrottled) {
    return;
  }

  isThrottled = true;

  setTimeout(async () => {
    const postsContainer = document.getElementById("posts");

    if (
      postsContainer.scrollTop + postsContainer.clientHeight >=
      postsContainer.scrollHeight
    ) {
      // Load more posts
      const user = Utils.getUser();
      if (category === "my-posts") {
        await getUserPosts(user.id); // Fetch user's created posts
      } else if (category === "liked-posts") {
        await getUserLikedPosts(user.id); // Fetch user's liked posts
      } else if (category === "disliked-posts") {
        await getUserDislikedPosts(user.id); // Fetch user's disliked posts
      } else {
        await getPostsByCategory(category); // Fetch posts by category
      }
    }
    isThrottled = false;
  }, 200);
}

export default class extends AbstractView {
  constructor(params) {
    super(params);
    this.setTitle("Forum");
  }

  async getHtml() {
    return `
      <div class="posts-box"><div id="posts" class="posts-grid">
        <!-- Content will be dynamically populated -->
      </div>
    `;
  }

  async init() {
    const postsGrid = document.querySelector(".posts-grid");
    postsGrid.textContent = "";
    // Initial load of posts
    await getPostsByCategory("General");

    // add scroll event listener
    const postsContainer = document.getElementById("posts");
    postsContainer?.addEventListener("scroll", handleScroll);

    const user = Utils.getUser();

    document.addEventListener("category-selected", async (e) => {
      const category_choosed = e.detail.category;

      offset = 0; 
      const postsDoc = document.querySelector(".posts-grid");
      postsDoc.textContent = ""; // Clear previous posts

      if (category_choosed === "my-posts") {
        await getUserPosts(user.id); // Fetch user's created posts
      } else if (category_choosed === "liked-posts") {
        await getUserLikedPosts(user.id); // Fetch user's liked posts
      } else if (category_choosed === "disliked-posts") {
        await getUserDislikedPosts(user.id); // Fetch user's disliked posts
      } else {
        await getPostsByCategory(category_choosed); // Fetch posts by category
      }

      category = category_choosed;
    });
  }
}
