class Timer {
	/**
	 *
	 * @param {HTMLElement} timerContainer
	 */
	constructor(timerContainer) {
		console.log(timerContainer, ' timer container');
		this.timerContainer = timerContainer;
		this.setupState(this.timerContainer);
		this.bindEvents();

		if (this.isRunning) {
			this.startUpdatingDisplay();
		}
	}

	/**
	 *
	 * @param {HTMLElement} dataContainer
	 */
	setupState(dataContainer) {
		this.timerDisplay = this.timerContainer.querySelector('.js-timer-display');
		this.startBtn = this.timerContainer.querySelector('.js-start-btn');
		this.pauseBtn = this.timerContainer.querySelector('.js-pause-btn');
		this.form = this.timerContainer.querySelector('.js-form');
		this.deleteBtn = this.timerContainer.querySelector('.js-delete-btn');
		console.log(this, ' this');
		this.titleInput = this.timerContainer.querySelector('.js-timer-title')

		this.id = +dataContainer.dataset['id'];
		this.seconds = +dataContainer.dataset['duration'];
		this.interval = null;
		const runningSince = dataContainer.dataset['runningSince'];
		const pausedAt = dataContainer.dataset['paused_at'];
		console.log(runningSince, pausedAt, 'hey');
		this.isRunning = !!runningSince && !pausedAt;
		console.log(this.isRunning, ' is running', this.id);
	}

	bindEvents() {
		this.startBtn.addEventListener('click', () => this.start());
		this.pauseBtn.addEventListener('click', () => this.pause());
		this.form.addEventListener('submit', (e) => this.handleSubmit(e));
		this.deleteBtn.addEventListener('click', (e) => this.delete(e));
		this.titleInput.addEventListener('blur', (e) => this.handleTimerTitleChange(e))
	}

	async handleTimerTitleChange(e) {
		const target = e.target;
		const newTitle = target.value;
		fetch(window.updateTimerTitle, {
			body:JSON.stringify({newTitle}),
			method:"POST",
			headers: {
				'Content-Type': 'application/json',
			}
		})
	}

	async start() {
		console.log('start');
		if (!this.isRunning) {
			this.isRunning = true;
		}
		console.log('timer id', this.timerContainer.dataset.id)
		const resp = await fetch(window.startTimer, {
			body: JSON.stringify({ timer_id: this.timerContainer.dataset.id }),
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
		});
		const data = await resp.json();
		console.log(data,' DARTA');
		this.startUpdatingDisplay();
	}

	pause() {
		if (!this.isRunning) return;

		this.isRunning = false;
		this.startBtn.disabled = false;
		this.pauseBtn.disabled = true;
		clearInterval(this.interval);

		const userId = this.timerContainer.dataset.userId;
		const timer_id = this.timerContainer.dataset.id;

		if (!userId || !timer_id) {
			alert('Не удалось определить ID юзера или таймера');
			return;
		}

		const pause_time = Date.now();
		fetch(window.pauseTimer, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ userId, timer_id, pause_time }),
		});
	}

	startUpdatingDisplay() {
		console.log('start updating display');
		this.startBtn.disabled = true;
		this.pauseBtn.disabled = false;

		this.interval = setInterval(() => {
			console.log(this.seconds,' seconds')
			this.seconds++;
			this.updateDisplay();
		}, 1000);
	}

	updateDisplay() {
		const hours = Math.floor(this.seconds / 3600);
		const remainingSeconds = this.seconds % 3600;
		const minutes = Math.floor(remainingSeconds / 60);
		const seconds = remainingSeconds % 60;
		const formattedTime = `${hours.toString().padStart(2, '0')}:${minutes
			.toString()
			.padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
		this.timerDisplay.textContent = formattedTime;
	}

	reset() {
		this.seconds = 0;
		this.updateDisplay();
		this.form.reset();
	}

	delete() {
		this.manager.removeTimer(this.id);
	}

	cleanup() {
		if (this.isRunning) {
			this.pause();
		}
	}
}
