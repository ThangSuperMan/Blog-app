const templateEditProfileName = `
    <form action="/edit_profile" method="post">
      <div class="text-field">
        <h4 for="profile_new_name">Type your new profile name</h4>
        <div class="profile-template-content">
          <input type="hidden" name="edit_profile_name" value="edit_profile_name" /> 
          <input type="text" name="profile_new_name"placeholder="New profile name" id="profile_new_name" required />
          <input type="password" name="password" placeholder="Password" id="password" required  />
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
`

const templateEditPassword = `
    <form action="/edit_profile" method="post">
      <div class="text-field">
        <h4 for="profile_new_name">Type your new password</h4>
        <div class="profile-template-content">
          <input type="hidden" name="edit_password" value="edit_password" /> 
          <input type="password" name="current_password" placeholder="Current password" id="current_password" required />
          <input type="password" name="new_password" placeholder="Password" id="password" required />
          <input type="password" name="confirm_new_password" placeholder="Confirm password" id="confirm_password" required />
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
`

function toArray(x) {
  var arr = []
  for (var i = 0; i < x.length; i++) {
    arr.push(x[i])
  }

  return arr
}

function handleClickEditButton(event) {
  let currentIdEditButton = event.target.getAttribute("id")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateBackground = document.querySelector(".profile-template-background")

  if (currentIdEditButton == "edit_profile_name") {
    profileTemplate.innerHTML = templateEditProfileName
  } else if (currentIdEditButton == "edit_password") {
    profileTemplate.innerHTML = templateEditPassword
  }

  profileTemplate.classList.add("active")
  profileTemplateBackground.classList.add("active")
}

// function handleClickAddBlog(event) {
//   console.log("handleClickAddBlog")
//   event.preventDefault()
// }
//
window.addEventListener("DOMContentLoaded", () => {
  // const addBlogBtn = document.querySelector("#add_blog")
  const editProfileNameBtn = document.querySelector("#edit_profile_name")
  const editPasswordBtn = document.querySelector("#edit_password")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateBackground = document.querySelector(".profile-template-background")
  const edits = document.getElementsByClassName("edit")
  editButtons = toArray(edits)

  // addBlogBtn.addEventListener("click", handleClickAddBlog)
  editProfileNameBtn.addEventListener('click', handleClickEditButton)
  editPasswordBtn.addEventListener('click', handleClickEditButton)

  profileTemplateBackground.addEventListener("mousedown", (e) => {

    if (e.target == profileTemplateBackground) {
      profileTemplate.classList.remove("active")
      profileTemplateBackground.classList.remove("active")
    }
  })
})

