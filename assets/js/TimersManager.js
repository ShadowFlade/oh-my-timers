class TimerManager {
	/**
	 *
	 * @param {HTMLElement} timersContainer
	 */
	constructor(timersContainer) {
		this.timers = [];
		this.timersContainer = timersContainer;
		console.log(this.timersContainer)

		this.addTimerBtn = document.getElementById('addTimerBtn');
		this.timersContainer = document.querySelector('.js-timers');

		this.bindEvents();
		this.initTimers();
	}

	bindEvents() {
		this.addTimerBtn.addEventListener('click', () => this.addTimer());
	}

	initTimers() {
		Array.from(this.timersContainer.children).forEach((timerHtml) => {
			const timer = new Timer(timerHtml);
			timer.manager = this
			console.log(timerHtml,' timer html');
		});
	}

	async addTimer() {
		const cookie = new Cookie();
		const userId = cookie.get('user_id');
		if(!userId) {
			alert('Не обнаружено id пользователя')
		}
		
		let newTimerHtml = '';
		const resp = await fetch(window.createTimer, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(
				{ userId }
			),
		})
		newTimerHtml = await resp.text()

		const htmlElement = this.createTimerFromRawHtml(newTimerHtml);
		this.timersContainer.appendChild(htmlElement);
		const timer = new Timer(htmlElement);
		this.timers.push(timer);
		timer.bindEvents();
	}

	/**
	 * 
	 * @param {string} rawHtml 
	 * @returns {HTMLElement}
	 */
	createTimerFromRawHtml(rawHtml) {
		const range = document.createRange();
		const htmlElement = range.createContextualFragment(rawHtml.trim()).firstChild;
		return htmlElement;
	}

	/**
	 * 
	 * @param {number} timerId 
	 */
	removeTimer(timerId) {
		const timerIndex = this.timers.findIndex((t) => t.id === timerId);
		if (timerIndex !== -1) {
			this.timers[timerIndex].cleanup();
			this.timers.splice(timerIndex, 1);

			const timerElement = document.getElementById(`timer-${timerId}`);
			if (timerElement) {
				timerElement.remove();
			}
		}
	}
}
