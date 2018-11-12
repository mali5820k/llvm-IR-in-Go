package asm

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/mewkiz/pkg/diffutil"
	"github.com/mewkiz/pkg/osutil"
)

// words specifies whether to colour words in diff output.
var words bool

func init() {
	flag.BoolVar(&words, "words", false, "colour words in diff output")
	flag.Parse()
}

func TestParseFile(t *testing.T) {
	golden := []struct {
		path string
	}{
		{path: "testdata/hexfloat.ll"},
		{path: "testdata/inst_aggregate.ll"},
		{path: "testdata/inst_binary.ll"},
		{path: "testdata/inst_bitwise.ll"},
		{path: "testdata/inst_conversion.ll"},
		{path: "testdata/inst_memory.ll"},
		{path: "testdata/inst_other.ll"},
		{path: "testdata/inst_vector.ll"},
		{path: "testdata/terminator.ll"},

		// LLVM test/Features.
		{path: "testdata/Feature/OperandBundles/adce.ll"},
		{path: "testdata/Feature/OperandBundles/basic-aa-argmemonly.ll"},
		{path: "testdata/Feature/OperandBundles/dse.ll"},
		{path: "testdata/Feature/OperandBundles/early-cse.ll"},
		{path: "testdata/Feature/OperandBundles/function-attrs.ll"},
		{path: "testdata/Feature/OperandBundles/inliner-conservative.ll"},
		{path: "testdata/Feature/OperandBundles/merge-func.ll"},
		{path: "testdata/Feature/OperandBundles/pr26510.ll"},
		{path: "testdata/Feature/OperandBundles/special-state.ll"},
		//{path: "testdata/Feature/alias2.ll"}, // TODO: fix grammar. syntax error at line 12
		//{path: "testdata/Feature/aliases.ll"}, // TODO: fix grammar. syntax error at line 29
		//{path: "testdata/Feature/alignment.ll"}, // TODO: fix grammar. syntax error at line 7
		{path: "testdata/Feature/attributes.ll"},
		{path: "testdata/Feature/basictest.ll"},
		{path: "testdata/Feature/callingconventions.ll"},
		{path: "testdata/Feature/calltest.ll"},
		{path: "testdata/Feature/casttest.ll"},
		{path: "testdata/Feature/cfgstructures.ll"},
		{path: "testdata/Feature/cold.ll"},
		{path: "testdata/Feature/comdat.ll"},
		//{path: "testdata/Feature/constexpr.ll"}, // TODO: re-enable when signed hex integer literals are supported.
		{path: "testdata/Feature/constpointer.ll"},
		{path: "testdata/Feature/const_pv.ll"},
		{path: "testdata/Feature/elf-linker-options.ll"},
		{path: "testdata/Feature/escaped_label.ll"},
		{path: "testdata/Feature/exception.ll"},
		{path: "testdata/Feature/float.ll"},
		{path: "testdata/Feature/fold-fpcast.ll"},
		{path: "testdata/Feature/forwardreftest.ll"},
		{path: "testdata/Feature/fp-intrinsics.ll"},
		{path: "testdata/Feature/global_pv.ll"},
		//{path: "testdata/Feature/globalredefinition3.ll"}, // TODO: figure out how to test. should report error "redefinition of global '@B'"
		{path: "testdata/Feature/global_section.ll"},
		{path: "testdata/Feature/globalvars.ll"},
		{path: "testdata/Feature/indirectcall2.ll"},
		{path: "testdata/Feature/indirectcall.ll"},
		{path: "testdata/Feature/inlineasm.ll"},
		{path: "testdata/Feature/instructions.ll"},
		{path: "testdata/Feature/intrinsic-noduplicate.ll"},
		{path: "testdata/Feature/intrinsics.ll"},
		{path: "testdata/Feature/load_module.ll"},
		{path: "testdata/Feature/md_on_instruction.ll"},
		{path: "testdata/Feature/memorymarkers.ll"},
		{path: "testdata/Feature/metadata.ll"},
		{path: "testdata/Feature/minsize_attr.ll"},
		{path: "testdata/Feature/NamedMDNode2.ll"},
		{path: "testdata/Feature/NamedMDNode.ll"},
		{path: "testdata/Feature/newcasts.ll"},
		{path: "testdata/Feature/optnone.ll"},
		{path: "testdata/Feature/optnone-llc.ll"},
		{path: "testdata/Feature/optnone-opt.ll"},
		{path: "testdata/Feature/packed.ll"},
		{path: "testdata/Feature/packed_struct.ll"},
		{path: "testdata/Feature/paramattrs.ll"},
		{path: "testdata/Feature/ppcld.ll"},
		{path: "testdata/Feature/prefixdata.ll"},
		{path: "testdata/Feature/prologuedata.ll"},
		{path: "testdata/Feature/properties.ll"},
		{path: "testdata/Feature/prototype.ll"},
		{path: "testdata/Feature/recursivetype.ll"},
		{path: "testdata/Feature/seh-nounwind.ll"},
		{path: "testdata/Feature/simplecalltest.ll"},
		{path: "testdata/Feature/smallest.ll"},
		{path: "testdata/Feature/small.ll"},
		{path: "testdata/Feature/sparcld.ll"},
		{path: "testdata/Feature/strip_names.ll"},
		//{path: "testdata/Feature/terminators.ll"}, // TODO: fix grammar. syntax error at line 35
		{path: "testdata/Feature/testalloca.ll"},
		{path: "testdata/Feature/testconstants.ll"},
		{path: "testdata/Feature/testlogical.ll"},
		//{path: "testdata/Feature/testtype.ll"}, // TODO: fix nil pointer dereference
		{path: "testdata/Feature/testvarargs.ll"},
		{path: "testdata/Feature/undefined.ll"},
		{path: "testdata/Feature/unreachable.ll"},
		{path: "testdata/Feature/varargs.ll"},
		{path: "testdata/Feature/varargs_new.ll"},
		{path: "testdata/Feature/vector-cast-constant-exprs.ll"},
		{path: "testdata/Feature/weak_constant.ll"},
		{path: "testdata/Feature/weirdnames.ll"},
		{path: "testdata/Feature/x86ld.ll"},

		// LLVM test/DebugInfo/Generic.
		//{path: "testdata/DebugInfo/Generic/2009-10-16-Phi.ll"},
		//{path: "testdata/DebugInfo/Generic/2009-11-03-InsertExtractValue.ll"},
		//{path: "testdata/DebugInfo/Generic/2009-11-05-DeadGlobalVariable.ll"},
		//{path: "testdata/DebugInfo/Generic/2009-11-06-NamelessGlobalVariable.ll"},
		//{path: "testdata/DebugInfo/Generic/2009-11-10-CurrentFn.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-01-05-DbgScope.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-03-12-llc-crash.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-03-19-DbgDeclare.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-03-24-MemberFn.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-04-06-NestedFnDbgInfo.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-04-19-FramePtr.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-05-03-DisableFramePtr.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-05-03-OriginDIE.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-05-10-MultipleCU.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-06-29-InlinedFnLocalVar.ll"},
		//{path: "testdata/DebugInfo/Generic/2010-10-01-crash.ll"},
		//{path: "testdata/DebugInfo/Generic/accel-table-hash-collisions.ll"},
		//{path: "testdata/DebugInfo/Generic/array.ll"},
		//{path: "testdata/DebugInfo/Generic/block-asan.ll"},
		//{path: "testdata/DebugInfo/Generic/bug_null_debuginfo.ll"},
		{path: "testdata/DebugInfo/Generic/constant-pointers.ll"},
		//{path: "testdata/DebugInfo/Generic/containing-type-extension.ll"},
		//{path: "testdata/DebugInfo/Generic/cross-cu-inlining.ll"},
		//{path: "testdata/DebugInfo/Generic/cross-cu-linkonce-distinct.ll"},
		//{path: "testdata/DebugInfo/Generic/cross-cu-linkonce.ll"},
		//{path: "testdata/DebugInfo/Generic/cu-range-hole.ll"},
		//{path: "testdata/DebugInfo/Generic/cu-ranges.ll"},
		//{path: "testdata/DebugInfo/Generic/dbg-at-specficiation.ll"},
		//{path: "testdata/DebugInfo/Generic/dead-argument-order.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-info-always-inline.ll"},
		{path: "testdata/DebugInfo/Generic/debug-info-enum.ll"},
		//{path: "testdata/DebugInfo/Generic/debuginfofinder-forward-declaration.ll"},
		//{path: "testdata/DebugInfo/Generic/debuginfofinder-imported-global-variable.ll"},
		//{path: "testdata/DebugInfo/Generic/debuginfofinder-inlined-cu.ll"},
		//{path: "testdata/DebugInfo/Generic/debuginfofinder-multiple-cu.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-info-qualifiers.ll"},
		{path: "testdata/DebugInfo/Generic/debug-label-mi.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-label-opt.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-empty-cu.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-empty-name.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-hash-collisions.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-index-type.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-linkage-name.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/Generic/debug-names-many-cu.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-name-collisions.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-one-cu.ll"},
		//{path: "testdata/DebugInfo/Generic/debug-names-two-cu.ll"},
		//{path: "testdata/DebugInfo/Generic/def-line.ll"},
		//{path: "testdata/DebugInfo/Generic/discriminated-union.ll"},
		//{path: "testdata/DebugInfo/Generic/discriminator.ll"},
		//{path: "testdata/DebugInfo/Generic/disubrange_vla.ll"},
		//{path: "testdata/DebugInfo/Generic/disubrange_vla_no_dbgvalue.ll"},
		//{path: "testdata/DebugInfo/Generic/dwarf-public-names.ll"},
		//{path: "testdata/DebugInfo/Generic/empty.ll"},
		//{path: "testdata/DebugInfo/Generic/enum.ll"},
		//{path: "testdata/DebugInfo/Generic/enum-types.ll"},
		//{path: "testdata/DebugInfo/Generic/extended-loc-directive.ll"},
		//{path: "testdata/DebugInfo/Generic/global.ll"},
		//{path: "testdata/DebugInfo/Generic/global-sra-array.ll"},
		//{path: "testdata/DebugInfo/Generic/global-sra-single-member.ll"},
		//{path: "testdata/DebugInfo/Generic/global-sra-struct.ll"},
		//{path: "testdata/DebugInfo/Generic/gmlt_profiling.ll"},
		//{path: "testdata/DebugInfo/Generic/gvn.ll"},
		//{path: "testdata/DebugInfo/Generic/imported-name-inlined.ll"},
		//{path: "testdata/DebugInfo/Generic/incorrect-variable-debugloc1.ll"},
		//{path: "testdata/DebugInfo/Generic/incorrect-variable-debugloc.ll"},
		//{path: "testdata/DebugInfo/Generic/indvar-discriminator.ll"},
		//{path: "testdata/DebugInfo/Generic/inheritance.ll"},
		//{path: "testdata/DebugInfo/Generic/inlined-arguments.ll"},
		//{path: "testdata/DebugInfo/Generic/inline-debug-info.ll"},
		//{path: "testdata/DebugInfo/Generic/inline-debug-info-multiret.ll"},
		//{path: "testdata/DebugInfo/Generic/inline-debug-loc.ll"},
		//{path: "testdata/DebugInfo/Generic/inlined-strings.ll"},
		//{path: "testdata/DebugInfo/Generic/inlined-vars.ll"},
		//{path: "testdata/DebugInfo/Generic/inline-no-debug-info.ll"},
		//{path: "testdata/DebugInfo/Generic/inline-scopes.ll"},
		//{path: "testdata/DebugInfo/Generic/instcombine-phi.ll"},
		{path: "testdata/DebugInfo/Generic/invalid.ll"},
		//{path: "testdata/DebugInfo/Generic/licm-hoist-debug-loc.ll"},
		//{path: "testdata/DebugInfo/Generic/linear-dbg-value.ll"},
		//{path: "testdata/DebugInfo/Generic/linkage-name-abstract.ll"},
		//{path: "testdata/DebugInfo/Generic/location-verifier.ll"},
		//{path: "testdata/DebugInfo/Generic/lto-comp-dir.ll"},
		//{path: "testdata/DebugInfo/Generic/mainsubprogram.ll"},
		//{path: "testdata/DebugInfo/Generic/member-order.ll"},
		//{path: "testdata/DebugInfo/Generic/member-pointers.ll"},
		//{path: "testdata/DebugInfo/Generic/missing-abstract-variable.ll"},
		//{path: "testdata/DebugInfo/Generic/multiline.ll"},
		//{path: "testdata/DebugInfo/Generic/namespace_function_definition.ll"},
		//{path: "testdata/DebugInfo/Generic/namespace_inline_function_definition.ll"},
		//{path: "testdata/DebugInfo/Generic/namespace.ll"},
		//{path: "testdata/DebugInfo/Generic/noscopes.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/Generic/pass-by-value.ll"},
		//{path: "testdata/DebugInfo/Generic/piece-verifier.ll"},
		//{path: "testdata/DebugInfo/Generic/PR20038.ll"},
		//{path: "testdata/DebugInfo/Generic/PR37395.ll"},
		//{path: "testdata/DebugInfo/Generic/ptrsize.ll"},
		//{path: "testdata/DebugInfo/Generic/recursive_inlining.ll"}, // TODO: fix grammar. syntax error at line 118
		//{path: "testdata/DebugInfo/Generic/restrict.ll"},
		//{path: "testdata/DebugInfo/Generic/simplifycfg_sink_last_inst.ll"},
		//{path: "testdata/DebugInfo/Generic/skeletoncu.ll"},
		//{path: "testdata/DebugInfo/Generic/sroa-larger.ll"},
		//{path: "testdata/DebugInfo/Generic/sroa-samesize.ll"},
		//{path: "testdata/DebugInfo/Generic/store-tail-merge.ll"},
		//{path: "testdata/DebugInfo/Generic/string-offsets-form.ll"},
		//{path: "testdata/DebugInfo/Generic/sugared-constants.ll"},
		//{path: "testdata/DebugInfo/Generic/sunk-compare.ll"},
		{path: "testdata/DebugInfo/Generic/template-recursive-void.ll"},
		//{path: "testdata/DebugInfo/Generic/thrownTypes.ll"},
		//{path: "testdata/DebugInfo/Generic/tu-composite.ll"},
		//{path: "testdata/DebugInfo/Generic/tu-member-pointer.ll"},
		//{path: "testdata/DebugInfo/Generic/two-cus-from-same-file.ll"},
		//{path: "testdata/DebugInfo/Generic/typedef.ll"},
		//{path: "testdata/DebugInfo/Generic/unconditional-branch.ll"},
		//{path: "testdata/DebugInfo/Generic/univariant-discriminated-union.ll"},
		//{path: "testdata/DebugInfo/Generic/varargs.ll"},
		//{path: "testdata/DebugInfo/Generic/version.ll"},
		//{path: "testdata/DebugInfo/Generic/virtual-index.ll"},
		//{path: "testdata/DebugInfo/Generic/volatile-alloca.ll"},

		// LLVM test/DebugInfo/X86.
		//{path: "testdata/DebugInfo/X86/2010-04-13-PubType.ll"},
		//{path: "testdata/DebugInfo/X86/2011-09-26-GlobalVarContext.ll"},
		//{path: "testdata/DebugInfo/X86/2011-12-16-BadStructRef.ll"},
		//{path: "testdata/DebugInfo/X86/abstract_origin.ll"},
		//{path: "testdata/DebugInfo/X86/accel-tables-dwarf5.ll"},
		//{path: "testdata/DebugInfo/X86/accel-tables.ll"},
		//{path: "testdata/DebugInfo/X86/align_c11.ll"},
		//{path: "testdata/DebugInfo/X86/align_cpp11.ll"},
		//{path: "testdata/DebugInfo/X86/aligned_stack_var.ll"},
		//{path: "testdata/DebugInfo/X86/align_objc.ll"},
		//{path: "testdata/DebugInfo/X86/arange-and-stub.ll"},
		//{path: "testdata/DebugInfo/X86/arange.ll"},
		//{path: "testdata/DebugInfo/X86/arguments.ll"},
		//{path: "testdata/DebugInfo/X86/array2.ll"},
		//{path: "testdata/DebugInfo/X86/array.ll"},
		//{path: "testdata/DebugInfo/X86/atomic-c11-dwarf-4.ll"},
		//{path: "testdata/DebugInfo/X86/atomic-c11-dwarf-5.ll"},
		//{path: "testdata/DebugInfo/X86/bbjoin.ll"},
		//{path: "testdata/DebugInfo/X86/bitcast-di.ll"},
		//{path: "testdata/DebugInfo/X86/bitfields-dwarf4.ll"},
		//{path: "testdata/DebugInfo/X86/bitfields.ll"},
		//{path: "testdata/DebugInfo/X86/block-capture.ll"},
		//{path: "testdata/DebugInfo/X86/byvalstruct.ll"},
		{path: "testdata/DebugInfo/X86/clang-module.ll"},
		//{path: "testdata/DebugInfo/X86/clone-module-2.ll"},
		//{path: "testdata/DebugInfo/X86/clone-module.ll"},
		//{path: "testdata/DebugInfo/X86/coff_debug_info_type.ll"},
		//{path: "testdata/DebugInfo/X86/coff_relative_names.ll"},
		//{path: "testdata/DebugInfo/X86/concrete_out_of_line.ll"},
		//{path: "testdata/DebugInfo/X86/constant-aggregate.ll"},
		//{path: "testdata/DebugInfo/X86/constant-loclist.ll"},
		//{path: "testdata/DebugInfo/X86/containing-type-extension-rust.ll"},
		//{path: "testdata/DebugInfo/X86/c-type-units.ll"},
		//{path: "testdata/DebugInfo/X86/cu-ranges.ll"},
		//{path: "testdata/DebugInfo/X86/cu-ranges-odr.ll"},
		//{path: "testdata/DebugInfo/X86/data_member_location.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-abstract-vars-g-gmlt.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-addr-dse.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-addr.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-byval-parameter.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-const-int.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-const.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-declare-alloca.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-declare-arg.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-declare-inalloca.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-declare.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-file-name.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-i128-const.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-merge-loc-entry.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-prolog-end.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-subrange.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-const-byref.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-dag-combine.ll"},
		//{path: "testdata/DebugInfo/X86/dbg_value_direct.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-frame-index.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-g-gmlt.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-inlined-parameter.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-isel.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-location.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-range.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-regmask-clobber.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-terminator.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-value-transfer-order.ll"},
		//{path: "testdata/DebugInfo/X86/dbg-vector-size.ll"},
		//{path: "testdata/DebugInfo/X86/debug_addr.ll"},
		//{path: "testdata/DebugInfo/X86/debug_and_nodebug_CUs.ll"},
		//{path: "testdata/DebugInfo/X86/debug-dead-local-var.ll"},
		//{path: "testdata/DebugInfo/X86/debug_frame.ll"},
		//{path: "testdata/DebugInfo/X86/debugger-tune.ll"},
		//{path: "testdata/DebugInfo/X86/debug-info-access.ll"},
		//{path: "testdata/DebugInfo/X86/debug-info-block-captured-self.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/debug-info-blocks.ll"},
		//{path: "testdata/DebugInfo/X86/debug-info-packed-struct.ll"},
		//{path: "testdata/DebugInfo/X86/debug-info-producer-with-flags.ll"},
		//{path: "testdata/DebugInfo/X86/debug-info-static-member.ll"},
		//{path: "testdata/DebugInfo/X86/debug-loc-asan.ll"},
		//{path: "testdata/DebugInfo/X86/debug-loc-frame.ll"},
		//{path: "testdata/DebugInfo/X86/debug-loc-offset.ll"},
		//{path: "testdata/DebugInfo/X86/debug-macro.ll"},
		//{path: "testdata/DebugInfo/X86/debug-names-split-dwarf.ll"},
		{path: "testdata/DebugInfo/X86/debug-ranges-offset.ll"},
		//{path: "testdata/DebugInfo/X86/decl-derived-member.ll"},
		//{path: "testdata/DebugInfo/X86/default-subrange-array.ll"},
		//{path: "testdata/DebugInfo/X86/deleted-bit-piece.ll"},
		{path: "testdata/DebugInfo/X86/DIModuleContext.ll"},
		{path: "testdata/DebugInfo/X86/DIModule.ll"},
		//{path: "testdata/DebugInfo/X86/discriminator2.ll"},
		//{path: "testdata/DebugInfo/X86/discriminator3.ll"},
		//{path: "testdata/DebugInfo/X86/discriminator.ll"},
		//{path: "testdata/DebugInfo/X86/dllimport.ll"},
		//{path: "testdata/DebugInfo/X86/double-declare.ll"},
		//{path: "testdata/DebugInfo/X86/dwarf-aranges.ll"},
		//{path: "testdata/DebugInfo/X86/dwarf-aranges-no-dwarf-labels.ll"},
		//{path: "testdata/DebugInfo/X86/dwarf-linkage-names.ll"},
		//{path: "testdata/DebugInfo/X86/dwarf-no-source-loc.ll"},
		//{path: "testdata/DebugInfo/X86/dwarf-public-names.ll"},
		//{path: "testdata/DebugInfo/X86/dwarf-pubnames-split.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_byte_size.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_calling-convention.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_linkage_name.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_location-reference.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_object_pointer.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_specification.ll"},
		//{path: "testdata/DebugInfo/X86/DW_AT_stmt_list_sec_offset.ll"},
		//{path: "testdata/DebugInfo/X86/dw_op_minus_direct.ll"},
		//{path: "testdata/DebugInfo/X86/dw_op_minus.ll"}, // TODO: add support for TLSModel initialexec.
		//{path: "testdata/DebugInfo/X86/DW_TAG_friend.ll"},
		//{path: "testdata/DebugInfo/X86/earlydup-crash.ll"},
		//{path: "testdata/DebugInfo/X86/elf-names.ll"},
		//{path: "testdata/DebugInfo/X86/empty-and-one-elem-array.ll"},
		//{path: "testdata/DebugInfo/X86/empty-array.ll"},
		//{path: "testdata/DebugInfo/X86/empty.ll"},
		//{path: "testdata/DebugInfo/X86/empty_macinfo.ll"},
		//{path: "testdata/DebugInfo/X86/ending-run.ll"},
		//{path: "testdata/DebugInfo/X86/enum-class.ll"},
		//{path: "testdata/DebugInfo/X86/enum-fwd-decl.ll"},
		//{path: "testdata/DebugInfo/X86/fi-expr.ll"},
		//{path: "testdata/DebugInfo/X86/fi-piece.ll"},
		//{path: "testdata/DebugInfo/X86/fission-cu.ll"},
		//{path: "testdata/DebugInfo/X86/fission-hash.ll"},
		//{path: "testdata/DebugInfo/X86/fission-inline.ll"},
		//{path: "testdata/DebugInfo/X86/fission-no-inlining.ll"},
		//{path: "testdata/DebugInfo/X86/fission-ranges.ll"},
		//{path: "testdata/DebugInfo/X86/float_const.ll"},
		//{path: "testdata/DebugInfo/X86/float_const_loclist.ll"},
		//{path: "testdata/DebugInfo/X86/formal_parameter.ll"},
		//{path: "testdata/DebugInfo/X86/fragment-offset-order.ll"},
		//{path: "testdata/DebugInfo/X86/FrameIndexExprs.ll"},
		//{path: "testdata/DebugInfo/X86/frame-register.ll"},
		//{path: "testdata/DebugInfo/X86/generate-odr-hash.ll"},
		//{path: "testdata/DebugInfo/X86/ghost-sdnode-dbgvalues.ll"},
		//{path: "testdata/DebugInfo/X86/global-expression.ll"},
		//{path: "testdata/DebugInfo/X86/global-sra-fp80-array.ll"},
		//{path: "testdata/DebugInfo/X86/global-sra-fp80-struct.ll"},
		//{path: "testdata/DebugInfo/X86/gnu-public-names-empty.ll"}, // TODO: fix grammar. syntax error at line 20
		//{path: "testdata/DebugInfo/X86/gnu-public-names-gmlt.ll"},
		//{path: "testdata/DebugInfo/X86/gnu-public-names.ll"},
		//{path: "testdata/DebugInfo/X86/gnu-public-names-multiple-cus.ll"},
		//{path: "testdata/DebugInfo/X86/gnu-public-names-tu.ll"},
		//{path: "testdata/DebugInfo/X86/header.ll"},
		//{path: "testdata/DebugInfo/X86/inline-asm-locs.ll"},
		//{path: "testdata/DebugInfo/X86/InlinedFnLocalVar.ll"},
		//{path: "testdata/DebugInfo/X86/inlined-formal-parameter.ll"},
		//{path: "testdata/DebugInfo/X86/inlined-indirect-value.ll"},
		//{path: "testdata/DebugInfo/X86/inline-member-function.ll"},
		//{path: "testdata/DebugInfo/X86/inline-namespace.ll"},
		//{path: "testdata/DebugInfo/X86/inline-seldag-test.ll"},
		//{path: "testdata/DebugInfo/X86/instcombine-instrinsics.ll"},
		//{path: "testdata/DebugInfo/X86/invalid-prologue-end.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/isel-cse-line.ll"},
		//{path: "testdata/DebugInfo/X86/lexical-block-file-inline.ll"},
		//{path: "testdata/DebugInfo/X86/lexical_block.ll"},
		//{path: "testdata/DebugInfo/X86/line-info.ll"},
		//{path: "testdata/DebugInfo/X86/linkage-name.ll"},
		//{path: "testdata/DebugInfo/X86/live-debug-values.ll"},
		//{path: "testdata/DebugInfo/X86/live-debug-variables.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/live-debug-vars-keep-undef.ll"},
		//{path: "testdata/DebugInfo/X86/low-pc-cu.ll"},
		//{path: "testdata/DebugInfo/X86/mem2reg_fp80.ll"},
		//{path: "testdata/DebugInfo/X86/memberfnptr.ll"},
		//{path: "testdata/DebugInfo/X86/mi-print.ll"},
		//{path: "testdata/DebugInfo/X86/misched-dbg-value.ll"},
		//{path: "testdata/DebugInfo/X86/missing-file-line.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/mixed-nodebug-cu.ll"},
		//{path: "testdata/DebugInfo/X86/multiple-aranges.ll"},
		//{path: "testdata/DebugInfo/X86/multiple-at-const-val.ll"},
		//{path: "testdata/DebugInfo/X86/nodebug.ll"},
		//{path: "testdata/DebugInfo/X86/no_debug_ranges.ll"},
		//{path: "testdata/DebugInfo/X86/nodebug_with_debug_loc.ll"},
		//{path: "testdata/DebugInfo/X86/nondefault-subrange-array.ll"},
		//{path: "testdata/DebugInfo/X86/nophysreg.ll"}, // TODO: add support for Dereferenceable
		//{path: "testdata/DebugInfo/X86/no-public-sections.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/noreturn_c11.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/noreturn_cpp11.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/noreturn_objc.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/objc-fwd-decl.ll"},
		//{path: "testdata/DebugInfo/X86/objc-property-void.ll"}, // TODO: add support for DIObjCProperty
		//{path: "testdata/DebugInfo/X86/op_deref.ll"},
		//{path: "testdata/DebugInfo/X86/parameters.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/partial-constant.ll"},
		//{path: "testdata/DebugInfo/X86/pieces-1.ll"},
		//{path: "testdata/DebugInfo/X86/pieces-2.ll"},
		//{path: "testdata/DebugInfo/X86/pieces-3.ll"},
		//{path: "testdata/DebugInfo/X86/pieces-4.ll"},
		//{path: "testdata/DebugInfo/X86/pointer-type-size.ll"},
		//{path: "testdata/DebugInfo/X86/pr11300.ll"},
		//{path: "testdata/DebugInfo/X86/pr12831.ll"},
		//{path: "testdata/DebugInfo/X86/pr13303.ll"},
		//{path: "testdata/DebugInfo/X86/pr19307.ll"},
		//{path: "testdata/DebugInfo/X86/PR26148.ll"},
		//{path: "testdata/DebugInfo/X86/pr28270.ll"},
		//{path: "testdata/DebugInfo/X86/pr34545.ll"},
		//{path: "testdata/DebugInfo/X86/PR37234.ll"},
		//{path: "testdata/DebugInfo/X86/processes-relocations.ll"},
		//{path: "testdata/DebugInfo/X86/prologue-stack.ll"},
		//{path: "testdata/DebugInfo/X86/range_reloc_base.ll"},
		//{path: "testdata/DebugInfo/X86/range_reloc.ll"},
		//{path: "testdata/DebugInfo/X86/ref_addr_relocation.ll"},
		//{path: "testdata/DebugInfo/X86/reference-argument.ll"},
		//{path: "testdata/DebugInfo/X86/rematerialize.ll"},
		//{path: "testdata/DebugInfo/X86/rnglists_base_attr.ll"},
		//{path: "testdata/DebugInfo/X86/rnglists_curanges.ll"},
		//{path: "testdata/DebugInfo/X86/rvalue-ref.ll"},
		//{path: "testdata/DebugInfo/X86/safestack-byval.ll"}, // TODO: add support for TLSModel initialexec.
		{path: "testdata/DebugInfo/X86/sdag-combine.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-dangling-dbgvalue.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-dbgvalue-phi-use-1.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-dbgvalue-phi-use-2.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-dbgvalue-phi-use-3.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-dbgvalue-phi-use-4.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-salvage-add.ll"},
		//{path: "testdata/DebugInfo/X86/sdagsplit-1.ll"},
		//{path: "testdata/DebugInfo/X86/sdag-split-arg.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/sections_as_references.ll"},
		//{path: "testdata/DebugInfo/X86/single-dbg_value.ll"},
		//{path: "testdata/DebugInfo/X86/single-fi.ll"},
		//{path: "testdata/DebugInfo/X86/spill-indirect-nrvo.ll"},
		//{path: "testdata/DebugInfo/X86/spill-nontrivial-param.ll"},
		//{path: "testdata/DebugInfo/X86/spill-nospill.ll"},
		//{path: "testdata/DebugInfo/X86/split-dwarf-cross-unit-reference.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		//{path: "testdata/DebugInfo/X86/split-dwarf-multiple-cu-hash.ll"},
		//{path: "testdata/DebugInfo/X86/split-dwarf-omit-empty.ll"},
		//{path: "testdata/DebugInfo/X86/split-global.ll"},
		//{path: "testdata/DebugInfo/X86/sret.ll"},
		//{path: "testdata/DebugInfo/X86/sroasplit-1.ll"},
		//{path: "testdata/DebugInfo/X86/sroasplit-2.ll"},
		//{path: "testdata/DebugInfo/X86/sroasplit-3.ll"},
		//{path: "testdata/DebugInfo/X86/sroasplit-4.ll"},
		//{path: "testdata/DebugInfo/X86/sroasplit-5.ll"},
		//{path: "testdata/DebugInfo/X86/sroasplit-dbg-declare.ll"},
		//{path: "testdata/DebugInfo/X86/stack-args.ll"},
		//{path: "testdata/DebugInfo/X86/stack-value-dwarf2.ll"}, // TODO: add support for Dereferenceable.
		//{path: "testdata/DebugInfo/X86/stack-value-dwarf4.ll"},
		//{path: "testdata/DebugInfo/X86/stack-value-piece.ll"},
		//{path: "testdata/DebugInfo/X86/static_member_array.ll"},
		//{path: "testdata/DebugInfo/X86/stmt-list.ll"},
		//{path: "testdata/DebugInfo/X86/stmt-list-multiple-compile-units.ll"},
		//{path: "testdata/DebugInfo/X86/string-offsets-multiple-cus.ll"},
		//{path: "testdata/DebugInfo/X86/string-offsets-table.ll"},
		//{path: "testdata/DebugInfo/X86/stringpool.ll"},
		//{path: "testdata/DebugInfo/X86/strip-broken-debuginfo.ll"},
		//{path: "testdata/DebugInfo/X86/struct-loc.ll"},
		//{path: "testdata/DebugInfo/X86/subrange-type.ll"},
		//{path: "testdata/DebugInfo/X86/subregisters.ll"},
		//{path: "testdata/DebugInfo/X86/subreg.ll"},
		//{path: "testdata/DebugInfo/X86/tail-merge.ll"},
		//{path: "testdata/DebugInfo/X86/template.ll"},
		//{path: "testdata/DebugInfo/X86/this-stack_value.ll"},
		//{path: "testdata/DebugInfo/X86/tls.ll"},
		//{path: "testdata/DebugInfo/X86/type_units_with_addresses.ll"},
		//{path: "testdata/DebugInfo/X86/unattached-global.ll"},
		//{path: "testdata/DebugInfo/X86/union-const.ll"},
		//{path: "testdata/DebugInfo/X86/union-template.ll"},
		//{path: "testdata/DebugInfo/X86/vector.ll"},
		//{path: "testdata/DebugInfo/X86/vla-dependencies.ll"},
		//{path: "testdata/DebugInfo/X86/vla-global.ll"},
		//{path: "testdata/DebugInfo/X86/vla.ll"},
		//{path: "testdata/DebugInfo/X86/vla-multi.ll"},
		//{path: "testdata/DebugInfo/X86/void-typedef.ll"},
		//{path: "testdata/DebugInfo/X86/xray-split-dwarf-interaction.ll"},
		//{path: "testdata/DebugInfo/X86/zextload.ll"},

		// LLVM test/DebugInfo.
		{path: "testdata/DebugInfo/check-debugify-preserves-analyses.ll"},
		{path: "testdata/DebugInfo/cross-cu-scope.ll"},
		{path: "testdata/DebugInfo/debugify-bogus-dbg-value.ll"},
		{path: "testdata/DebugInfo/debugify-each.ll"},
		{path: "testdata/DebugInfo/debugify-export.ll"},
		{path: "testdata/DebugInfo/debugify.ll"},
		{path: "testdata/DebugInfo/debugify-report-missing-locs-only.ll"},
		//{path: "testdata/DebugInfo/debuglineinfo-path.ll"}, // TODO: figure out how to handle AttrGroupID with missing AttrGroupDef
		{path: "testdata/DebugInfo/dwo.ll"},
		{path: "testdata/DebugInfo/macro_link.ll"},
		{path: "testdata/DebugInfo/omit-empty.ll"},
		{path: "testdata/DebugInfo/pr34186.ll"},
		{path: "testdata/DebugInfo/pr34672.ll"},
		{path: "testdata/DebugInfo/skeletoncu.ll"},
		{path: "testdata/DebugInfo/strip-DIGlobalVariable.ll"},
		{path: "testdata/DebugInfo/strip-loop-metadata.ll"},
		{path: "testdata/DebugInfo/strip-module-flags.ll"},
		//{path: "testdata/DebugInfo/unrolled-loop-remainder.ll"}, // TODO: figure out how to handle duplicate (but distinct) AttrGroupDef

		// Coreutils.
		/*
			{path: "testdata/coreutils/[.ll"},
			{path: "testdata/coreutils/b2sum.ll"},
			{path: "testdata/coreutils/base32.ll"},
			{path: "testdata/coreutils/base64.ll"},
			{path: "testdata/coreutils/basename.ll"},
			{path: "testdata/coreutils/cat.ll"},
			{path: "testdata/coreutils/chcon.ll"},
			{path: "testdata/coreutils/chgrp.ll"},
			{path: "testdata/coreutils/chmod.ll"},
			{path: "testdata/coreutils/chown.ll"},
			{path: "testdata/coreutils/chroot.ll"},
			{path: "testdata/coreutils/cksum.ll"},
			{path: "testdata/coreutils/comm.ll"},
			{path: "testdata/coreutils/cp.ll"},
			{path: "testdata/coreutils/csplit.ll"},
			{path: "testdata/coreutils/cut.ll"},
			{path: "testdata/coreutils/date.ll"},
			{path: "testdata/coreutils/dd.ll"},
			{path: "testdata/coreutils/df.ll"},
			{path: "testdata/coreutils/dir.ll"},
			{path: "testdata/coreutils/dircolors.ll"},
			{path: "testdata/coreutils/dirname.ll"},
			{path: "testdata/coreutils/du.ll"},
			{path: "testdata/coreutils/echo.ll"},
			{path: "testdata/coreutils/env.ll"},
			{path: "testdata/coreutils/expand.ll"},
			{path: "testdata/coreutils/expr.ll"},
			{path: "testdata/coreutils/factor.ll"},
			{path: "testdata/coreutils/false.ll"},
			{path: "testdata/coreutils/fmt.ll"},
			{path: "testdata/coreutils/fold.ll"},
			{path: "testdata/coreutils/getlimits.ll"},
			{path: "testdata/coreutils/ginstall.ll"},
			{path: "testdata/coreutils/groups.ll"},
			{path: "testdata/coreutils/head.ll"},
			{path: "testdata/coreutils/hostid.ll"},
			{path: "testdata/coreutils/id.ll"},
			{path: "testdata/coreutils/join.ll"},
			{path: "testdata/coreutils/kill.ll"},
			{path: "testdata/coreutils/link.ll"},
			{path: "testdata/coreutils/ln.ll"},
			{path: "testdata/coreutils/logname.ll"},
			{path: "testdata/coreutils/ls.ll"},
			{path: "testdata/coreutils/make-prime-list.ll"},
			{path: "testdata/coreutils/md5sum.ll"},
			{path: "testdata/coreutils/mkdir.ll"},
			{path: "testdata/coreutils/mkfifo.ll"},
			{path: "testdata/coreutils/mknod.ll"},
			{path: "testdata/coreutils/mktemp.ll"},
			{path: "testdata/coreutils/mv.ll"},
			{path: "testdata/coreutils/nice.ll"},
			{path: "testdata/coreutils/nl.ll"},
			{path: "testdata/coreutils/nohup.ll"},
			{path: "testdata/coreutils/nproc.ll"},
			{path: "testdata/coreutils/numfmt.ll"},
			{path: "testdata/coreutils/od.ll"},
			{path: "testdata/coreutils/paste.ll"},
			{path: "testdata/coreutils/pathchk.ll"},
			{path: "testdata/coreutils/pinky.ll"},
			{path: "testdata/coreutils/pr.ll"},
			{path: "testdata/coreutils/printenv.ll"},
			{path: "testdata/coreutils/printf.ll"},
			{path: "testdata/coreutils/ptx.ll"},
			{path: "testdata/coreutils/pwd.ll"},
			{path: "testdata/coreutils/readlink.ll"},
			{path: "testdata/coreutils/realpath.ll"},
			{path: "testdata/coreutils/rm.ll"},
			{path: "testdata/coreutils/rmdir.ll"},
			{path: "testdata/coreutils/runcon.ll"},
			{path: "testdata/coreutils/seq.ll"},
			{path: "testdata/coreutils/sha1sum.ll"},
			{path: "testdata/coreutils/sha224sum.ll"},
			{path: "testdata/coreutils/sha256sum.ll"},
			{path: "testdata/coreutils/sha384sum.ll"},
			{path: "testdata/coreutils/sha512sum.ll"},
			{path: "testdata/coreutils/shred.ll"},
			{path: "testdata/coreutils/shuf.ll"},
			{path: "testdata/coreutils/sleep.ll"},
			{path: "testdata/coreutils/sort.ll"},
			{path: "testdata/coreutils/split.ll"},
			{path: "testdata/coreutils/stat.ll"},
			{path: "testdata/coreutils/stdbuf.ll"},
			{path: "testdata/coreutils/stty.ll"},
			{path: "testdata/coreutils/sum.ll"},
			{path: "testdata/coreutils/sync.ll"},
			{path: "testdata/coreutils/tac.ll"},
			{path: "testdata/coreutils/tail.ll"},
			{path: "testdata/coreutils/tee.ll"},
			{path: "testdata/coreutils/test.ll"},
			{path: "testdata/coreutils/timeout.ll"},
			{path: "testdata/coreutils/touch.ll"},
			{path: "testdata/coreutils/tr.ll"},
			{path: "testdata/coreutils/true.ll"},
			{path: "testdata/coreutils/truncate.ll"},
			{path: "testdata/coreutils/tsort.ll"},
			{path: "testdata/coreutils/tty.ll"},
			{path: "testdata/coreutils/uname.ll"},
			{path: "testdata/coreutils/unexpand.ll"},
			{path: "testdata/coreutils/uniq.ll"},
			{path: "testdata/coreutils/unlink.ll"},
			{path: "testdata/coreutils/uptime.ll"},
			{path: "testdata/coreutils/users.ll"},
			{path: "testdata/coreutils/vdir.ll"},
			{path: "testdata/coreutils/wc.ll"},
			{path: "testdata/coreutils/who.ll"},
			{path: "testdata/coreutils/whoami.ll"},
			{path: "testdata/coreutils/yes.ll"},
		*/

		// SQLite.
		//{path: "testdata/sqlite/shell.ll"},
	}
	for _, g := range golden {
		log.Printf("=== [ %s ] ===", g.path)
		m, err := ParseFile(g.path)
		if err != nil {
			t.Errorf("unable to parse %q into AST; %+v", g.path, err)
			continue
		}
		path := g.path
		if osutil.Exists(g.path + ".golden") {
			path = g.path + ".golden"
		}
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("unable to read %q; %+v", path, err)
			continue
		}
		want := string(buf)
		got := m.String()
		if want != got {
			if err := diffutil.Diff(want, got, words, filepath.Base(path)); err != nil {
				panic(err)
			}
			t.Errorf("module mismatch %q; expected `%s`, got `%s`", path, want, got)
			continue
		}
	}
}
