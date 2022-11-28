//Test Javascript - not yet fully working

const form = document.getElementById('createEventForm');

form.addEventListener('submit', (event) => {
    // stop form submission
    document.getElementById("test").innerHTML = "JAVASCRIPT WAS HERE! WORKS!!";
});

function myFunction() {
	document.getElementById("test").innerHTML = "JAVASCRIPT WAS HERE! WORKS!!";
}