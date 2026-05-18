document.addEventListener("DOMContentLoaded", () => {

  const slices = document.querySelectorAll(".slice");
  const counter = document.getElementById("counter");
  const status = document.getElementById("status");

  let isDown = false;
  let mode = "add";
  let state = {};

  // --- EXISTING CODE FOR SLICE SELECTION ---
  slices.forEach((el, idx) => {
    el.addEventListener("mousedown", () => {
      isDown = true;
      mode = el.classList.contains("selected") ? "remove" : "add";
      apply(el, idx);
      updateCounter();
    });

    el.addEventListener("mouseover", () => {
      if (!isDown) return;
      apply(el, idx);
      updateCounter();
    });
  });

  document.addEventListener("mouseup", () => { isDown = false; });

  const saveBtn = document.getElementById("saveBtn");
  if (saveBtn) {
      saveBtn.addEventListener("click", () => {
          const selected = [];
          document.querySelectorAll(".slice.selected").forEach(el => {
              selected.push({
                  day: el.dataset.day,
                  hour: el.dataset.hour,
                  slice: Number(el.dataset.slice)
              });
          });

          const days = ["Mon","Tue","Wed","Thu","Fri"];
          let dailyCounts = {};
          days.forEach(day => { dailyCounts[day] = 0; });

          Object.keys(state).forEach(key => {
              if (!state[key]) return;

              const [day] = key.split("-");
              dailyCounts[day]++;
          });

          for (let day of days) {
              if (dailyCounts[day] < 18) {
                  status.textContent = `Day ${day} has too few slices (min 18)`;
                  return;
              } else if (dailyCounts[day] > 54) {
                  status.textContent = `Day ${day} has too many slices (max 54)`;
                  return;
              }
          }

          fetch("/save-selection", {
              method: "POST",
              headers: {"Content-Type": "application/json"},
              body: JSON.stringify({ selected })
          })
          .then(r => r.text())
          .then(t => status.textContent = t)
          .catch(() => status.textContent = "Error saving");
      });
  } else {
      console.warn("saveBtn not found in the DOM");
  }
  

  // Approve/Deny buttons (guarded)
  const approveBtn = document.getElementById("approveBtn");
  if (approveBtn) {
      approveBtn.addEventListener("click", () => scheduleApproveDeny("approve"));
  }

  const denyBtn = document.getElementById("denyBtn");
  if (denyBtn) {
      denyBtn.addEventListener("click", () => scheduleApproveDeny("deny"));
  }

  // Return to menu (always exists)
  document.getElementById("retMenu").addEventListener("click", () => {
      window.location.href = "/menu";
  });

  // --- helper functions ---
  function apply(el, idx) {
      if (mode === "add") {
          el.classList.add("selected");
          state[idx] = true;
      } else {
          el.classList.remove("selected");
          state[idx] = false;
      }
  }

  function updateCounter() {
    const selectedCells = document.querySelectorAll(".slice.selected");
    const count = selectedCells.length;

    // Weekly total (hours and minutes)
    const hrs = Math.floor(count / 6);
    const mins = (count % 6) * 10;
    counter.textContent = "Selected: " + hrs + "h " + mins + "m";

    const days = ["Mon","Tue","Wed","Thu","Fri"];
    let dailyCounts = {};
   
    days.forEach(day => { dailyCounts[day] = 0; });
   
    selectedCells.forEach(el => {
        dailyCounts[el.dataset.day]++;
    });
    
    days.forEach(day => {
        const span = document.getElementById(`counter-${day}`);
        if (!span) return;
        span.textContent = `${day}: ${dailyCounts[day]}`;
        // Highlight if outside limits
        if (dailyCounts[day] < 18 || dailyCounts[day] > 54) {
            span.style.color = "red";
        } else {
            span.style.color = "green";
        }
    });
}

});