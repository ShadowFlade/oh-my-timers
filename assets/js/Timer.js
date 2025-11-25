class Timer {
	/**
	 *
	 * @param {HTMLElement} timerContainer
	 */
	constructor(timerContainer) {
		this.timerContainer = timerContainer;
		this.setupState(this.timerContainer);
		this.initTime()
		this.bindEvents();

		if (this.isRunning) {
			this.startUpdatingDisplay();
		}
		this.cssClasses = {
			timerCircleRunning:"timer-circle--active"
		}
	}

	/**
	 *
	 * @param {HTMLElement} dataContainer
	 */
	setupState(dataContainer) {
		this.timerCircle = this.timerContainer.querySelector('.js-timer-circle')
		this.timerDisplay = this.timerContainer.querySelector('.js-timer-display');
		this.startBtn = this.timerContainer.querySelector('.js-start-btn');
		this.pauseBtn = this.timerContainer.querySelector('.js-pause-btn');
		this.stopBtn = this.timerContainer.querySelector('.js-stop-btn');
		this.form = this.timerContainer.querySelector('.js-form');
		this.deleteBtn = this.timerContainer.querySelector('.js-delete-btn');
		this.titleInput = this.timerContainer.querySelector('.js-timer-title')
		this.refreshButton = this.timerContainer.querySelector('.js-refresh-btn')

		this.id = +dataContainer.dataset['id'];
		this.seconds = +dataContainer.dataset['duration'];
		this.updatingDisplayInterval = null;
		const runningSince = dataContainer.dataset['runningSince'];
		const pausedAt = dataContainer.dataset['paused_at'];
		this.isRunning = !!runningSince && !pausedAt;
	}

	bindEvents() {
		this.startBtn.addEventListener('click', () => this.start());
		this.pauseBtn.addEventListener('click', () => this.pause());
		this.stopBtn.addEventListener('click',() => this.stop())
		this.refreshButton.addEventListener('click', () => this.refresh())
		this.form.addEventListener('submit', (e) => this.handleSubmit(e));
		this.deleteBtn.addEventListener('click', (e) => this.delete(e));
		this.titleInput.addEventListener('blur', (e) => this.handleTimerTitleChange(e))
	}

	async handleTimerTitleChange(e) {
		const target = e.target;
		const newTitle = target.value;
		const timerId = this.timerContainer.dataset.id;
		if (!timerId) {
			return;
		}
		fetch(window.updateTimerTitle, {
			body:JSON.stringify({newTitle, id: timerId}),
			method:"POST",
			headers: {
				'Content-Type': 'application/json',
			}
		})
	}

	async start() {
		if (!this.isRunning) {
			this.isRunning = true;
		}
		const resp = await fetch(window.startTimer, {
			body: JSON.stringify({
				timer_id: this.timerContainer.dataset.id,
				start_time: Date.now() / 1000,
			}),
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
		});
		const data = await resp.json();
		this.startUpdatingDisplay();
		this.updateCssClasses();
	}

	initTime() {
		const runningSince = this.timerContainer.dataset.runningSince;
		if (!runningSince) {
			return; //TODO[quality]:make frontend logger (send to backend)
		}
		const runningSinceDate = new Date(runningSince);
		const runningSinceTime = runningSinceDate.getTime()
		const now = Date.now();
		console.log({now,runningSinceDate,runningSinceTime})
		const seconds = Math.round((now - runningSinceTime) / 1000);
		this.seconds = seconds;
		console.log(this.timerContainer.dataset.runningSince, 'running since', seconds, 'seconds');
	}

	pause() {
		if (!this.isRunning) return;

		this.isRunning = false;
		this.startBtn.disabled = false;
		this.pauseBtn.disabled = true;
	
		clearInterval(this.updatingDisplayInterval);

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

	/**
	 * Сбрасывает время и останавливает таймер
	 */
	async stop() {
		const stop_time = Date.now()
		this.seconds = 0;
		this.updateDisplay();
		clearInterval(this.updatingDisplayInterval);
		this.runn
		this.isRunning = false;
		this.startBtn.disabled = false;
		this.pauseBtn.disabled = true;
		const userId = this.timerContainer.dataset.userId;
		const timer_id = this.timerContainer.dataset.id;

		if (!userId || !timer_id) {
			alert('Не удалось определить ID юзера или таймера');
			return;
		}
		const resp = await fetch(window.stopTimer, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ userId, timer_id, stop_time }),
		});
		const data = await resp.json();
		this.updateCssClasses();
		console.log(data,' data');
	}



	startUpdatingDisplay() {
		this.startBtn.disabled = true;
		this.pauseBtn.disabled = false;
		this.stopBtn.disabled = false;
		

		this.updatingDisplayInterval = setInterval(() => {
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

	refresh() {
		clearInterval(this.updatingDisplayInterval)
		this.seconds = 0;
		this.updateDisplay()
		this.startUpdatingDisplay();
		const timerId = +this.timerContainer.dataset.id;
		fetch(window.refreshTimer, {
			body:JSON.stringify({timerId}),
			method:"POST",
			headers: {
				'Content-Type': 'application/json',
			}
		})
	}

	delete() {
		this.manager.removeTimer(this.id);
	}

	cleanup() {
		if (this.isRunning) {
			this.pause();
		}
	}

	updateCssClasses() {
		const circleRunningClass = this.cssClasses.timerCircleRunning
		
		if (
			this.isRunning 
			&& !this.timerCircle.classList.contains(circleRunningClass)
		) {
			this.timerCircle.classList.add(circleRunningClass)
		}
		console.log(this.timerCircle, ' timer circle', this.isRunning);

		if (
			!this.isRunning
			&& this.timerCircle.classList.contains(circleRunningClass)
		) {
			console.log('removing');
				this.timerCircle.classList.remove(circleRunningClass)
			}
	}
}
