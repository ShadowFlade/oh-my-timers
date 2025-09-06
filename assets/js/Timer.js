class Timer {
	constructor(id, manager) {
		this.id = id;
		this.manager = manager;
		this.seconds = 0;
		this.interval = null;
		this.isRunning = false;
	}
	
	bindEvents() {
		const container = document.querySelector(`js-timer-${this.id}`);
		
		this.timerDisplay = document.getElementById(`timerDisplay-${this.id}`);
		this.startBtn = container.querySelector('.start-btn');
		this.stopBtn = container.querySelector('.stop-btn');
		this.form = container.querySelector('.form');
		this.deleteBtn = container.querySelector('.delete-btn');
		
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