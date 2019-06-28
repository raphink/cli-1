package cmd

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

type templateListItemOutput struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreationDate string `json:"creation_date"`
	Zone         string `json:"zone"`
	DiskSize     string `json:"disk_size"`
}

type templateListOutput []templateListItemOutput

func (o *templateListOutput) toJSON()  { outputJSON(o) }
func (o *templateListOutput) toText()  { outputText(o) }
func (o *templateListOutput) toTable() { outputTable(o) }

func init() {
	templateListCmd.Flags().BoolP("community", "", false, "List community templates")
	templateListCmd.Flags().BoolP("featured", "", false, "List featured templates")
	templateListCmd.Flags().BoolP("mine", "", false, "List your templates")
	templateListCmd.Flags().StringP("zone", "", "", "Name of the zone (default: current account's default zone)")
	templateCmd.AddCommand(templateListCmd)
}

var templateListCmd = &cobra.Command{
	Use:   "list [keyword]",
	Short: "List all available templates",
	Long: fmt.Sprintf(`This command lists available Compute Instance templates. By default, returns "featured" templates.

Supported output template annotations: %s`,
		strings.Join(outputterTemplateAnnotations(&templateListOutput{}), ", ")),
	Aliases: gListAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		var templateFilter string

		zone, err := cmd.Flags().GetString("zone")
		if err != nil {
			return err
		}

		if zone == "" {
			zone = gCurrentAccount.DefaultZone
		}

		community, err := cmd.Flags().GetBool("community")
		if err != nil {
			return err
		}

		mine, err := cmd.Flags().GetBool("mine")
		if err != nil {
			return err
		}

		if community {
			templateFilter = "community"
		} else if mine {
			templateFilter = "self"
		} else {
			templateFilter = "featured"
		}

		return output(listTemplates(templateFilter, zone, args))
	},
}

func listTemplates(templateFilter, zone string, filters []string) (outputter, error) {
	zoneID, err := getZoneIDByName(zone)
	if err != nil {
		return nil, err
	}

	templates, err := findTemplates(zoneID, templateFilter, filters...)
	if err != nil {
		return nil, err
	}

	out := templateListOutput{}

	for _, template := range templates {
		out = append(out, templateListItemOutput{
			ID:           template.ID.String(),
			Name:         template.Name,
			DiskSize:     humanize.IBytes(uint64(template.Size)),
			CreationDate: template.Created,
			Zone:         template.ZoneName,
		})
	}

	return &out, nil
}
