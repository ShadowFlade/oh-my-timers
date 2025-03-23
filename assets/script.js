document.addEventListener("DOMContentLoaded", () => {
	const newTimerButton = document.querySelector(".js-new-timer__button");
	newTimerButton.addEventListener("click", async (e) => {
		e.preventDefault();
		const response = await fetch("createTimer");
		console.log(response,' respnse')
		const newTimerHTML = response.text
		/**
		 * @type HTMLElement 
		 */
		const blockToReplace = e.target.closest('.js-new-timer');
		const parent = blockToReplace.parentElement
		if(!blockToReplace.outerHTML) {
			return;
		}
		blockToReplace.outerHTML = newTimerHTML
		const newTimer = parent.querySelector(".js-timer")
		newTimer.addEventListener
	})


})