// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.purrl;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.purrl.ProviderArgs;
import com.pulumi.purrl.Utilities;
import javax.annotation.Nullable;

@ResourceType(type="pulumi:providers:purrl")
public class Provider extends com.pulumi.resources.ProviderResource {
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public Provider(String name) {
        this(name, ProviderArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public Provider(String name, @Nullable ProviderArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public Provider(String name, @Nullable ProviderArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("purrl", name, args == null ? ProviderArgs.Empty : args, makeResourceOptions(options, Codegen.empty()));
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
            .build();
        return com.pulumi.resources.CustomResourceOptions.merge(defaultOptions, options, id);
    }

}
