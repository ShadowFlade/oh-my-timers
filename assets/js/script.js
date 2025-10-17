document.addEventListener("DOMContentLoaded", () => {
	const timersContainer = document.querySelector('.js-timers')
	const timersManager = new TimerManager(timersContainer);
})


window.createTimer = '/createTimer';
window.createUser = '/createUser';
window.pauseTimer = '/pauseTimer';
window.stopTimer = '/stopTimer';
window.deleteTimer = '/deleteTimer';
window.startTimer = '/startTimer';

window.updateTimerTitle = '/updateTimerTitle';
window.home = '/';
