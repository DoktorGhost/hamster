const FULL_DASH_ARRAY = 283;
const WARNING_THRESHOLD = 10; // Оставьте как есть, в секундах
const ALERT_THRESHOLD = 5;   // Оставьте как есть, в секундах

const COLOR_CODES = {
    info: {
        color: "green"
    },
    warning: {
        color: "orange",
        threshold: WARNING_THRESHOLD
    },
    alert: {
        color: "red",
        threshold: ALERT_THRESHOLD
    }
};

let timePassed = 0;
let timerInterval = null;
let remainingPathColor = COLOR_CODES.info.color;

function setupTimerUI(timeLimitMinutes) {
    const timeLimitSeconds = timeLimitMinutes * 60;
    document.getElementById("app").innerHTML = `
    <div class="base-timer">
      <svg class="base-timer__svg" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
        <g class="base-timer__circle">
          <circle class="base-timer__path-elapsed" cx="50" cy="50" r="45"></circle>
          <path
            id="base-timer-path-remaining"
            stroke-dasharray="283"
            class="base-timer__path-remaining ${remainingPathColor}"
            d="
              M 50, 50
              m -45, 0
              a 45,45 0 1,0 90,0
              a 45,45 0 1,0 -90,0
            "
          ></path>
        </g>
      </svg>
      <span id="base-timer-label" class="base-timer__label">${formatTime(timeLimitSeconds)}</span>
    </div>
    `;
}

function onTimesUp() {
    clearInterval(timerInterval);
    // Можно добавить дополнительную обработку по окончании таймера
}

function startTimer(timeLimitMinutes) {
    const timeLimitSeconds = timeLimitMinutes * 60;
    timePassed = 0;
    timerInterval = setInterval(() => {
        timePassed += 1;
        const timeLeft = timeLimitSeconds - timePassed;
        document.getElementById("base-timer-label").innerHTML = formatTime(timeLeft);
        setCircleDasharray(timeLeft, timeLimitSeconds);
        setRemainingPathColor(timeLeft);

        if (timeLeft <= 0) {
            onTimesUp();
        }
    }, 1000);
}

function formatTime(time) {
    const minutes = Math.floor(time / 60);
    const seconds = time % 60;
    return `${minutes}:${seconds < 10 ? `0${seconds}` : seconds}`;
}

function setRemainingPathColor(timeLeft) {
    const { alert, warning, info } = COLOR_CODES;
    if (timeLeft <= alert.threshold) {
        document
            .getElementById("base-timer-path-remaining")
            .classList.remove(warning.color);
        document
            .getElementById("base-timer-path-remaining")
            .classList.add(alert.color);
    } else if (timeLeft <= warning.threshold) {
        document
            .getElementById("base-timer-path-remaining")
            .classList.remove(info.color);
        document
            .getElementById("base-timer-path-remaining")
            .classList.add(warning.color);
    }
}

function calculateTimeFraction(timeLeft, timeLimit) {
    const rawTimeFraction = timeLeft / timeLimit;
    return rawTimeFraction - (1 / timeLimit) * (1 - rawTimeFraction);
}

function setCircleDasharray(timeLeft, timeLimit) {
    const circleDasharray = `${(
        calculateTimeFraction(timeLeft, timeLimit) * FULL_DASH_ARRAY
    ).toFixed(0)} 283`;
    document
        .getElementById("base-timer-path-remaining")
        .setAttribute("stroke-dasharray", circleDasharray);
}
