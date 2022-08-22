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
  console.log(profileTemplate)

  if (currentIdEditButton == "edit_profile_name") {
    console.log("dit_profile_name")
    profileTemplate.innerHTML = templateEditProfileName
  } else if (currentIdEditButton == "edit_password") {
    console.log("edit_password")
    profileTemplate.innerHTML = templateEditPassword
  }

  profileTemplate.classList.add("active")
  profileTemplateBackground.classList.add("active")
  // profileTemplate.innerHTML = templateEditProfileName
}

const templateEditProfileName = `
    <form action="edit_profile" method="post">
      <div class="text-field">
        <h4 for="profile_new_name">Type your new profile name</h4>
        <div class="profile-template-content">
          <input type="text" name="profile_new_name"placeholder="Profile new name" id="profile_new_name" />
          <input type="text" name="password" placeholder="Password" id="password" />
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
`

const templateEditPassword = `
    <form action="/profile/update" method="post">
      <div class="text-field">
        <h4 for="profile_new_name">Type your new password</h4>
        <div class="profile-template-content">
          <input type="text" placeholder="Current password" id="current_password" />
          <input type="text" placeholder="Password" id="password" />
          <input type="text" placeholder="Confirm password" id="confirm_password" />
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
`

window.addEventListener("DOMContentLoaded", () => {
  console.log('hello')
  const editProfileNameBtn = document.querySelector("#edit_profile_name")
  const editPasswordBtn = document.querySelector("#edit_password")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateBackground = document.querySelector(".profile-template-background")
  const edits = document.getElementsByClassName("edit")
  editButtons = toArray(edits)

  editProfileNameBtn.addEventListener('click', handleClickEditButton)
  editPasswordBtn.addEventListener('click', handleClickEditButton)

  // editButtons.forEach(button => button.addEventListener('click', handleClickhandleClickEditButtonEditButton))
  profileTemplateBackground.addEventListener("mousedown", (e) => {

    if (e.target == profileTemplateBackground) {
      console.log("You just clicked the profilt template background")
      profileTemplate.classList.remove("active")
      profileTemplateBackground.classList.remove("active")
      // profileTemplate.innerHTML = "<h1>What's up</h1>"
    }
  })
})

