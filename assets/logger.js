const Logger = {
  data: () => {
    return {
      weight: 0.0,
      reps: 0,
      sets: [],
      selectedSet: null,
      increaseWeight() {
        this.weight = parseFloat(this.weight) + 2.5;
      },
      decreaseWeight() {
        this.weight = parseFloat(this.weight) - 2.5;
      },
      increaseReps() {
        this.reps = parseInt(this.reps, 10) + 1;
      },
      decreaseReps() {
        this.reps = parseInt(this.reps, 10) - 1;
      },
      onSubmit() {
        this.sets.push({
          weight: this.weight,
          reps: this.reps,
        });
      },
      selectSet(set) {
        if (set === this.selectedSet) {
          this.selectedSet = null;
        } else {
          this.selectedSet = set;
          this.weight = set.weight;
          this.reps = set.reps;
        }
      },
      clear() {
        this.weight = 0.0;
        this.reps = 0;
      },
      update() {
        const thisSet = this.sets.find(s => s === this.selectedSet);
        thisSet.weight = this.weight;
        thisSet.reps = this.reps;
      },
      deleteSet() {
        this.sets = this.sets.filter(set => set !== this.selectedSet);
        this.selectedSet = null;
      },
    }
  },
};
