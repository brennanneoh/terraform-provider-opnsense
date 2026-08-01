package interfaces

import (
	"github.com/browningluke/opnsense-go/pkg/interfaces"
	"github.com/browningluke/terraform-provider-opnsense/internal/tools"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// assignmentResourceModel describes the resource data model.
type assignmentResourceModel struct {
	Device      types.String `tfsdk:"device"`
	Description types.String `tfsdk:"description"`
	Lock        types.Bool   `tfsdk:"lock"`
	Id          types.String `tfsdk:"id"`
}

func assignmentResourceSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Assigns a physical (or virtual) device to a logical interface, such as `wan`, `lan`, or an optional `optN` interface. " +
			"OPNsense's REST API does not currently expose the interface's general configuration (IPv4/IPv6 addressing mode, enable, MTU, etc.) " +
			"— only this device assignment, description, and lock flag can be managed here.\n\n" +
			"To manage the existing `wan` or `lan` interface, import it first: `terraform import opnsense_interfaces_assignment.wan wan`. " +
			"Creating a new assignment adds an additional optional interface (`opt1`, `opt2`, ...); `terraform destroy` unassigns the interface.",

		Attributes: map[string]schema.Attribute{
			"device": schema.StringAttribute{
				MarkdownDescription: "The physical (or virtual) device to assign to this interface, e.g. `vtnet0`.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Optional:            true,
			},
			"lock": schema.BoolAttribute{
				MarkdownDescription: "Prevent removal of the interface. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Logical interface identifier (e.g. `wan`, `lan`, `opt1`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func assignmentDataSourceSchema() dschema.Schema {
	return dschema.Schema{
		MarkdownDescription: "Assigns a physical (or virtual) device to a logical interface, such as `wan`, `lan`, or an optional `optN` interface.",

		Attributes: map[string]dschema.Attribute{
			"id": dschema.StringAttribute{
				MarkdownDescription: "Logical interface identifier (e.g. `wan`, `lan`, `opt1`).",
				Required:            true,
			},
			"device": dschema.StringAttribute{
				MarkdownDescription: "The physical (or virtual) device assigned to this interface.",
				Computed:            true,
			},
			"description": dschema.StringAttribute{
				MarkdownDescription: "Optional description here for your reference (not parsed).",
				Computed:            true,
			},
			"lock": dschema.BoolAttribute{
				MarkdownDescription: "Whether removal of the interface is prevented.",
				Computed:            true,
			},
		},
	}
}

func convertAssignmentSchemaToStruct(d *assignmentResourceModel) (*interfaces.Assignment, error) {
	return &interfaces.Assignment{
		Device:      d.Device.ValueString(),
		Description: d.Description.ValueString(),
		Lock:        tools.BoolToString(d.Lock.ValueBool()),
	}, nil
}

func convertAssignmentStructToSchema(d *interfaces.Assignment) (*assignmentResourceModel, error) {
	return &assignmentResourceModel{
		Device:      types.StringValue(d.Device),
		Description: tools.StringOrNull(d.Description),
		Lock:        types.BoolValue(tools.StringToBool(d.Lock)),
	}, nil
}
