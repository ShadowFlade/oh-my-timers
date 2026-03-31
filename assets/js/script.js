document.addEventListener("DOMContentLoaded", () => {
    const timersContainer = document.querySelector('.js-root')
    const timerSectionContainer = timersContainer.querySelector('.js-sections')
    const timerSections = timerSectionContainer.querySelectorAll('.js-section')
    timerSections.forEach(timerSection => new TimerManager(timerSection))
})


window.createTimer = '/createTimer';
window.createUser = '/createUser';
window.pauseTimer = '/pauseTimer';
window.stopTimer = '/stopTimer';
window.refreshTimer = '/refreshTimer'
window.deleteTimer = '/deleteTimer';
window.startTimer = '/startTimer';
window.updateTimerColor = "/updateTimerColor"
window.updateTimerTitle = '/updateTimerTitle';
window.home = '/';
