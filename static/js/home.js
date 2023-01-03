async function getAuthorsOfBlogs() {
  console.log("getAuthorsOfBlogs");

  const response = await fetch("http://localhost:3000/", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ year: yearEl.value }),
  });
  const data = await response.json();
}

window.addEventListener("DOMContentLoaded", () => {
  console.log("Hello from home.js");
  getAuthorsOfBlogs();
});
