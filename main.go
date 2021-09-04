package main

import (
    "bufio"
    "fmt"
    "flag"
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"

    "tasks"
)

// Path
var json_file = "task_data.json"

// Flags
var add_task *bool
var update_task *bool
var delete_task *bool
var display_tasks *bool
var sort *bool
var task_id *int
var task_priority *string
var task_name *string
var task_due *string

func init() {
    // Initialize flags.
    add_task = flag.Bool("add", false, "Defines if a task should be added.")
    update_task = flag.Bool("update", false, "Defines if a task should be updated.")
    delete_task = flag.Bool("delete", false, "Defines if a task should be deleted.")
    display_tasks = flag.Bool("display", true, "Displays table of tasks to console.")
    sort = flag.Bool("sort", false, "Will resort based on priority and due date.")
    task_id = flag.Int("id", 0, "Id of task to use.")
    task_priority = flag.String("priority", "", "Lower the number the more important.")
    task_name = flag.String("name", "", "Name of task to use.")
    task_due = flag.String("due", "", "When is the task due. Format: YYYY-MM-DD HH:MM AM")
}

func date_flag_valid(flag_val string) bool {
    m := "^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9] [0-9][0-9]:[0-9][0-9] [A|P]M$"
    result, err := regexp.Match(m, []byte(flag_val))
    if err != nil {
        log.Fatal(err)
    }
    return result
}

func date_flag_to_time(due_flag string) time.Time {
    result, err := time.Parse("2006-01-02 3:04 PM", due_flag)
    if err != nil {
        log.Fatal(err)
    }
    return result
}

func main() {
    // Read flags provided by user.
    flag.Parse()

    // Create a TaskMaster object.
    task_master := tasks.NewTaskMaster()

    // If the file that holds the json data for TaskMaster doesn't exist,
    // create it. Otherwise load the data frm the json file.
    if _, err := os.Stat(json_file); os.IsNotExist(err) {
        task_master.SaveToJson(json_file)
    } else {
        task_master.LoadFromJson(json_file)
    }

    // Add a new task.
    if *add_task == true && *task_name != "" {
        new_task := tasks.NewTask(*task_name)
        new_task.Priority, _ = strconv.Atoi(*task_priority)
        if date_flag_valid(*task_due) {
            new_task.Due = date_flag_to_time(*task_due)
        }
        task_id := task_master.RegisterTask(new_task)
        fmt.Printf(`Task "%s" created with id "%d".`, *task_name, task_id)
        fmt.Println()
    }

    // Update a task based on provided id.
    if *update_task == true {
        if len(task_master.Tasks) < (*task_id + 1) {
            fmt.Println("Task id does not exist")
            return
        }

        tsk := &task_master.Tasks[*task_id]

        // Verify update with user.
        rdr := bufio.NewReader(os.Stdin)
        fmt.Println("Are you sure you want to update this task:")
        fmt.Println("   ", *task_id,"    "+tsk.Name)
        fmt.Println("Type \"Y\" to confirm.")
    	os.Stdout.Write([]byte("#: ")) //+"\n" <---- Do I need this?
    	raw_answer, err := rdr.ReadString('\n')
    	if err != nil {
            log.Fatal(err)
    	}
    	cleanup_input := strings.NewReplacer("\n", "")
    	answer := strings.TrimSpace(cleanup_input.Replace(raw_answer))

        if answer == "Y" {
            if *task_name != "" {
                tsk.Name = *task_name
            }
            if *task_priority != "" {
                tsk.Priority, _ = strconv.Atoi(*task_priority)
            }
            if date_flag_valid(*task_due) {
                tsk.Due = date_flag_to_time(*task_due)
            }
        }
    }

    // Delete a task based on the provided id.
    if *delete_task == true {
        if len(task_master.Tasks) < (*task_id + 1) {
            fmt.Println("Task id does not exist")
            return
        }
        task_name := task_master.Tasks[*task_id].Name

        // Verify delition with user.
        rdr := bufio.NewReader(os.Stdin)
        fmt.Println("Are you sure you want to delete this task:")
        fmt.Println("   ", *task_id,"    "+task_name)
        fmt.Println("Type \"Y\" to confirm.")
    	os.Stdout.Write([]byte("#: ")) //+"\n" <---- Do I need this?
    	raw_answer, err := rdr.ReadString('\n')
    	if err != nil {
            log.Fatal(err)
    	}
    	cleanup_input := strings.NewReplacer("\n", "")
    	answer := strings.TrimSpace(cleanup_input.Replace(raw_answer))

        if answer == "Y" {
            task_master.DeleteTaskId(*task_id)
            fmt.Printf(`Task "%d" deleted.`, *task_id)
            fmt.Println()
        }
    }

    // Display all tasks.
    if *sort == true {
        task_master.SortByPriority()
    }

    // Display all tasks.
    if *display_tasks == true {
        task_master.TableDisplayTo(os.Stdout)
    }

    // Save updated data to the json_file.
    task_master.SaveToJson(json_file)
}
