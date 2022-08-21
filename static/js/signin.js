console.log("sign in js")

function handleSubmitForm(e) {
  // e.preventDefault()
  console.log("handleSubmitForm")
}

window.addEventListener("DOMContentLoaded", () => {
  console.log('DOMContentLoade')
  const form = document.querySelector(".form")
  console.log(form)
  form.addEventListener('submit', handleSubmitForm)
})
