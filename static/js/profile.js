function toArray(x) {
  var arr = []
  for (var i = 0; i < x.length; i++) {
    arr.push(x[i])
  }

  return arr
}

function handleClickEditButton(event) {
  // console.log()
  let currentIdEditButton = event.target.getAttribute("id")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateBackground = document.querySelector(".profile-template-background")

  if (currentIdEditButton == "edit_profile_name") {
    console.log("haha the id is edit_profile_name") 
  } else if (currentIdEditButton == "edit_password_name") {
    console.log("haha the id is edit_password_name") 
  }

  profileTemplate.classList.add("active")
  profileTemplateBackground.classList.add("active")
  profileTemplate.innerHTML = templateEditProfileName 
}

const templateEditProfileName = `
    <form action="/profile/update" method="post">
      <div class="text-field">
        <label for="profile_new_name">Type your new profile name</label>
        <div class="profile-template-content">
          <input type="text" placeholder="Username" id="profile_new_name" />
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
`

const templateEditPassword = `
    <form action="/profile/update" method="post">
      <div class="text-field">
        <label for="profile_new_name">Type your new profile name</label>
        <div class="profile-template-content">
          <input type="text" placeholder="Username" id="profile_new_name" />
          <button type="submit">Submit</button>
        </div>
      </div>
    </form>
`

window.addEventListener("DOMContentLoaded", () => {
  const editProfileNameBtn = document.querySelector("#edit_profile_name")
  const editPasswordBtn = document.querySelector("#edit_password")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateBackground = document.querySelector(".profile-template-background")
  const editButtons = document.getElementsByClassName("edit")
  // let editButtons = toArray(edits)

  
  editProfileNameBtn.addEventListener('click', handleClickEditButton)
  editPasswordBtn.addEventListener('click', handleClickEditButton)

  editButtons.forEach(button => button.addEventListener('click', handleClickEditButton))
  profileTemplateBackground.addEventListener("mousedown", (e) => {

    if (e.target == profileTemplateBackground) {
      console.log("e.target == profile-template-background")
      profileTemplate.classList.remove("active")
      profileTemplateBackground.classList.remove("active")
      profileTemplate.innerHTML = templateEditProfileName 
    }
  })
})

