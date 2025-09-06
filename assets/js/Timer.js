class Timer {
	/**
	 * 
	 * @param {number} id 
	 * @param {HTMLElement} timerContainer 
	 */
	constructor(id, timerContainer) {
		this.timerContainer = timerContainer

	
	}

	/**
	 * 
	 * @param {HTMLElement} dataContainer 
	 */
	setupState(dataContainer){
		this.id = +dataContainer.dataset["id"];
		this.seconds = +dataContainer.dataset["duration"];;
		this.interval = null;
		this.isRunning = !!dataContainer.dataset["runningSince"];
	}
	
	bindEvents() {		
		this.timerDisplay = this.timerContainer.querySelector("js-timer-display");
		this.startBtn = this.timerContainer.querySelector('.js-start-btn');
		this.stopBtn = this.timerContainer.querySelector('.js-stop-btn');
		this.form = this.timerContainer.querySelector('.js-form');
		this.deleteBtn = this.timerContainer.querySelector('.js-delete-btn');
		
		this.startBtn.addEventListener('click', () => this.start());
		this.stopBtn.addEventListener('click', () => this.stop());
		this.form.addEventListener('submit', (e) => this.handleSubmit(e));
		this.deleteBtn.addEventListener('click', () => this.delete());
	}
	
	start() {
		if (!this.isRunning) {
			this.isRunning = true;
			this.startBtn.disabled = true;
			this.stopBtn.disabled = false;
			
			this.interval = setInterval(() => {
				this.seconds++;
				this.updateDisplay();
			}, 1000);
		}
	}
	
	stop() {
		if (this.isRunning) {
			this.isRunning = false;
			this.startBtn.disabled = false;
			this.stopBtn.disabled = true;
			
			clearInterval(this.interval);
		}
	}
	
	updateDisplay() {
		const minutes = Math.floor(this.seconds / 60);
		const remainingSeconds = this.seconds % 60;
		
		const formattedTime = `${minutes.toString().padStart(2, '0')}:${remainingSeconds.toString().padStart(2, '0')}`;
		this.timerDisplay.textContent = formattedTime;
	}
	
	handleSubmit(e) {
		e.preventDefault();
		
		const formData = new FormData(this.form);
		const title = formData.get('title');
		
		if (this.isRunning) {
			this.stop();
		}
		
		console.log('PATCH request would be sent with:', {
			timerId: this.id,
			title: title,
			duration: this.seconds,
			formattedTime: this.timerDisplay.textContent
		});
		
		alert(`Timer "${title}" saved!\nDuration: ${this.timerDisplay.textContent}`);
		
		this.reset();
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
			this.stop();
		}
	}
}