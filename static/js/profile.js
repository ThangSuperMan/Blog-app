function toArray(x) {
  var arr = []
  for (var i = 0; i < x.length; i++) {
    arr.push(x[i])
  }

  return arr
}

function handleClickButton() {
  console.log("hhandleClickButton")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateOverlay = document.querySelector(".profile-template-overlay")
  console.log(profileTemplateOverlay )
  profileTemplateOverlay.classList.add("active")
  profileTemplate.classList.add("active")
  profileTemplate.innerHTML = template
}

function handleMosueDownBodyTag() {
  console.log("handleMosueEnterBodyTag")
  const profileTemplate = document.querySelector(".profile-template")
  const profileTemplateOverlay = document.querySelector(".profile-template-overlay")
  profileTemplate .classList.remove("active")
  profileTemplateOverlay.classList.remove("active")
}

const template = `
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
  const edits = document.getElementsByClassName("edit")
  let editButtons = toArray(edits)

  editButtons.forEach(button => button.addEventListener('click', handleClickButton))
  document.body.addEventListener("mousedown", handleMosueDownBodyTag)
  
})

