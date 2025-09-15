class TimerManager {
	/**
	 *
	 * @param {HTMLElement} timersContainer
	 */
	constructor(timersContainer) {
		this.checkForExistingUserAndRegister();
		this.timers = [];
		this.timersContainer = timersContainer;
		console.log(this.timersContainer)

		this.addTimerBtn = document.querySelector('.js-new-timer__button');
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

	/**
	 * Мб потом вынести в другую сущность
	 */
	checkForExistingUserAndRegister() {
		const cookie = new Cookie();
		const userId = cookie.get('user_id_detected');
		console.log(userId, ' user id from cookie',document.cookie);
		if (!userId) {
			const superSecretPassword = prompt('Не обнаружено id пользователя. Введите супер секретный пароль от вашего пользователя');
			const user = new User();
			user.createUser(superSecretPassword);
			return;
		}
	}

	async addTimer() {
		let newTimerHtml = '';
		const resp = await fetch(window.createTimer, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
		})
		newTimerHtml = await resp.text()

		const htmlElement = this.createTimerFromRawHtml(newTimerHtml);
		console.log(htmlElement,' html element',htmlElement.firstChild,htmlElement.firstElementChild);
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
