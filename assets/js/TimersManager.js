class TimerManager {
	/**
	 *
	 * @param {HTMLElement} timersContainer
	 */
	constructor(timersContainer) {
		this.checkForExistingUserAndRegister();
		this.timers = [];
		this.timersContainer = timersContainer;
		console.log(this.timersContainer);

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
			this.timers.push(timer);
			timer.manager = this;
			console.log(timerHtml, ' timer html');
		});
	}

	/**
	 * Мб потом вынести в другую сущность
	 */
	async checkForExistingUserAndRegister() {
		const cookie = new Cookie();
		const userId = cookie.get('user_id_detected');
		console.log(userId, ' user id from cookie', document.cookie);

		if (!userId || this.checkForNewUserTrigger()) {
			const superSecretPassword = prompt(
				'Не обнаружено id пользователя. Введите супер секретный пароль от вашего пользователя'
			);
			const user = new User();
			const resp = await user.createUser(superSecretPassword);
			const data = await resp.json();

			return data.isSuccess;
		}
	}

	/**
	 * Проверяем есть ли скрытый инпут, который говорит о том, что нужно показать alert для ввода юзером пароля для создания нового юзера
	 * @returns boolean Нужно ли показывать алерт для создания нового юзера
	 */
	checkForNewUserTrigger() {
		const triggerHiddenInput = document.querySelector('.js-show-new-user-alert-trigger');

		return !!triggerHiddenInput;
	}

	async addTimer() {
		let newTimerHtml = '';
		const resp = await fetch(window.createTimer, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
		});
		newTimerHtml = await resp.text();

		const htmlElement = this.createTimerFromRawHtml(newTimerHtml);
		console.log(
			htmlElement,
			' html element',
			htmlElement.firstChild,
			htmlElement.firstElementChild
		);
		this.timersContainer.appendChild(htmlElement);
		const timer = new Timer(htmlElement);
		this.timers.push(timer);
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
	async removeTimer(timerId) {
		console.log(timerId, 'timer id');
		const timerIndex = this.timers.findIndex((t) => t.id === timerId);
		console.log(this.timers, 'TIMERS');
		console.log(timerIndex, 'timerIndex');

		if (timerIndex !== -1) {
			const timerElement = this.timers[timerIndex].timerContainer;
			if (!timerElement) {
				console.error('Не удалось найти таймер с id ' + timerId);
				console.info(this.timers);
			}
			this.timers[timerIndex].cleanup();
			this.timers.splice(timerIndex, 1);

			if (timerElement) {
				const resp = await fetch(window.deleteTimer, {
					body: JSON.stringify({ timer_id: timerId }),
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
				});
				const data = await resp.json();
				console.log(data,' DATA')
				if (data.IsSuccess) {
					timerElement.remove();
				} else {
					alert("Не удалось удалить таймер. Анлак.")
				}
			}
		}
	}
}
