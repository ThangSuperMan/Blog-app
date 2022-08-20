function toArray(x) {
  var arr = []
  for(var i = 0; i < x.length; i++) {
    arr.push(x[i])
  }

  return arr
}

window.addEventListener("DOMContentLoaded", () => {
  const edits = document.getElementsByClassName("edit")

  // var startTime = performance.now()
  // sayHi("Ahihi, to day is super good")
  // var endTime = performance.now()
  // console.log("startTime: ", startTime)
  // console.log("endTime: ", endTime)
  // console.log(`Call to doSomething took ${endTime - startTime} milliseconds`)
  
  var startTime = performance.now()
  let myArray = toArray(edits)
  var endTime = performance.now()
  console.log("myArray: ", myArray)
  console.log(`Call to doSomething took ${endTime - startTime} milliseconds`)

  var arr = [].slice.call(edits)
  console.log("myArrray when using the slice: ", arr)


})

