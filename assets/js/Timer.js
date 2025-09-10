class Timer {
	/**
	 *
	 * @param {HTMLElement} timerContainer
	 */
	constructor(timerContainer) {
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
		this.deleteBtn.addEventListener('click', () => this.delete());
	}

	start() {
		console.log('start');
		if (!this.isRunning) {
			this.isRunning = true;
		}
		this.startUpdatingDisplay();
	}

	pause() {
		if (!this.isRunning) return;
		
		this.isRunning = false;
		this.startBtn.disabled = false;
		this.pauseBtn.disabled = true;
		clearInterval(this.interval);

		const userId = this.timerContainer.dataset.userId
		const timer_id = this.timerContainer.dataset.id;
		
		if (!userId || !timer_id) {
			alert("Не удалось определить ID юзера или таймера")
			return;
		}
		
		fetch(window.pauseTimer, {
			method:"POST",
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(
				{ userId, timer_id}
			),
			
		})
	}

	startUpdatingDisplay() {
		this.startBtn.disabled = true;
		this.pauseBtn.disabled = false;

		this.interval = setInterval(() => {
			this.seconds++;
			this.updateDisplay();
		}, 1000);
	}

	updateDisplay() {
		const minutes = Math.floor(this.seconds / 60);
		const remainingSeconds = this.seconds % 60;

		const formattedTime = `${minutes.toString().padStart(2, '0')}:${remainingSeconds
			.toString()
			.padStart(2, '0')}`;
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
