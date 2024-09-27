package job

import (
	"camp-mgr/app/campmgr/internal/config"
	"context"
	"sync"
	"testing"
	"time"
)

// TestAddCampaignJob tests adding a scheduled job
func TestAddCampaignJob(t *testing.T) {

	var c config.Config
	scheduler := NewJobClient(context.Background(), c, nil, nil).(*DefaultJobClient)
	scheduler.Start()
	defer scheduler.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	now := time.Now().Unix()
	scheduler.AddCampaignJob(123456, now+1)

	// Use WaitGroup to confirm task execution
	scheduler.cron.AddFunc("* * * * * *", func() {
		wg.Done() // Mark as done after the job is executed
	})

	wg.Wait()

	if _, exists := scheduler.jobIDs[123456]; !exists {
		t.Errorf("Scheduled job was not added correctly")
	}
}

// TestUpdateCampaignJob tests updating a scheduled job
func TestUpdateCampaignJob(t *testing.T) {

	var c config.Config
	scheduler := NewJobClient(context.Background(), c, nil, nil).(*DefaultJobClient)
	scheduler.Start()
	defer scheduler.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	now := time.Now().Unix()
	scheduler.AddCampaignJob(123457, now+2)

	// Ensure the job is executed after being updated
	scheduler.UpdateCampaignJob(123457, now+1)

	if _, exists := scheduler.jobIDs[123457]; !exists {
		t.Errorf("Old job was not deleted after update")
	}

	scheduler.cron.AddFunc("* * * * * *", func() {
		wg.Done() // Mark as done after the new job is executed
	})

	wg.Wait()
}

// TestRemoveCampaignJob tests removing a scheduled job
func TestRemoveCampaignJob(t *testing.T) {
	var c config.Config
	scheduler := NewJobClient(context.Background(), c, nil, nil).(*DefaultJobClient)
	scheduler.Start()
	defer scheduler.Stop()

	now := time.Now().Unix()
	scheduler.AddCampaignJob(123458, now+2)

	scheduler.RemoveCampaignJob(123458)

	if _, exists := scheduler.jobIDs[123458]; exists {
		t.Errorf("Scheduled job was not removed correctly")
	}
}
