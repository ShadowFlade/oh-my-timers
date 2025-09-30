class TimerManager {
	/**
	 *
	 * @param {HTMLElement} timersContainer
	 */
	constructor(timersContainer) {
		this.checkForExistingUserAndRegister();
		this.timers = [];
		this.timersContainer = timersContainer;

		this.addTimerBtn = document.querySelector('.js-new-timer__button');
		this.timersContainer = document.querySelector('.js-timers');

		this.bindEvents();
		this.initTimers();
	}

	bindEvents() {
		this.addTimerBtn.addEventListener('click', () => this.addTimer());
	}

	initTimers() {
		let tabIndex = 1;
		Array.from(this.timersContainer.children).forEach((timerHtml) => {
			const timer = new Timer(timerHtml);
			timer.titleInput.tabIndex = tabIndex// i dont like this shit but doing this in go templates is tedious
			this.timers.push(timer);
			timer.manager = this;
			tabIndex++;
		});

		this.lastTabIndex = tabIndex
	}

	/**
	 * Мб потом вынести в другую сущность
	 */
	async checkForExistingUserAndRegister() {
		const cookie = new Cookie();
		const userId = cookie.get('user_id_detected');

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

		this.timersContainer.appendChild(htmlElement);
		htmlElement.tabIndex= this.lastTabIndex + 1;
		this.lastTabIndex = this.lastTabIndex + 1;
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
		const timerIndex = this.timers.findIndex((t) => t.id === timerId);

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

				if (data.IsSuccess) {
					timerElement.remove();
				} else {
					alert("Не удалось удалить таймер. Анлак.")
				}
			}
		}
	}
}
